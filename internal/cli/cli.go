package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/wohsj110/figma_cli/api"
	"github.com/wohsj110/figma_cli/internal/credential"
	"github.com/wohsj110/figma_cli/internal/figmaurl"
)

type app struct {
	stdin   io.Reader
	stdout  io.Writer
	stderr  io.Writer
	verbose bool
}

func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	a := &app{stdin: stdin, stdout: stdout, stderr: stderr}
	return a.run(ctx, args)
}

func (a *app) run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}
	args = a.parseGlobal(args)
	if len(args) == 0 {
		a.printHelp()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		a.printHelp()
		return nil
	case "init":
		return a.runInit(args[1:])
	case "me":
		return a.runMe(ctx, args[1:])
	case "file":
		return a.runFile(ctx, args[1:])
	case "node":
		return a.runNode(ctx, args[1:])
	case "image":
		return a.runImage(ctx, args[1:])
	case "comments":
		return a.runComments(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *app) parseGlobal(args []string) []string {
	out := make([]string, 0, len(args))
	for _, arg := range args {
		switch arg {
		case "--verbose", "-v":
			a.verbose = true
		default:
			out = append(out, arg)
		}
	}
	return out
}

func (a *app) runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	tokenStdin := fs.Bool("token-stdin", false, "read Figma token from stdin")
	fromEnv := fs.String("from-env", "", "read Figma token from environment variable")
	if err := fs.Parse(interspersed(args, map[string]bool{"from-env": true})); err != nil {
		return err
	}
	switch {
	case *tokenStdin && *fromEnv != "":
		return errors.New("--token-stdin and --from-env are mutually exclusive")
	case !*tokenStdin && *fromEnv == "":
		return errors.New("no token source: pass --token-stdin or --from-env VAR")
	}
	var raw string
	if *fromEnv != "" {
		v, ok := os.LookupEnv(*fromEnv)
		if !ok || strings.TrimSpace(v) == "" {
			return fmt.Errorf("environment variable %s is unset or empty", *fromEnv)
		}
		raw = v
	} else {
		b, err := io.ReadAll(a.stdin)
		if err != nil {
			return fmt.Errorf("read token from stdin: %w", err)
		}
		raw = string(b)
	}
	if err := credential.Store(raw); err != nil {
		return err
	}
	fmt.Fprintln(a.stdout, "Stored Figma token in OS keyring")
	return nil
}

func (a *app) runMe(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("me", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	raw := fs.Bool("raw", false, "print raw JSON")
	if err := fs.Parse(interspersed(args, nil)); err != nil {
		return err
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	me, body, err := client.Me(ctx)
	if err != nil {
		return err
	}
	if *raw {
		return writeRaw(a.stdout, body)
	}
	fmt.Fprintf(a.stdout, "ID: %s\nHandle: %s\nEmail: %s\n", me.ID, me.Handle, me.Email)
	return nil
}

func (a *app) runFile(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("file command required: get")
	}
	switch args[0] {
	case "get":
		return a.runFileGet(ctx, args[1:])
	default:
		return fmt.Errorf("unknown file command %q", args[0])
	}
}

func (a *app) runFileGet(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("file get", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	depth := fs.Int("depth", 1, "document depth to fetch")
	raw := fs.Bool("raw", false, "print raw JSON")
	if err := fs.Parse(interspersed(args, map[string]bool{"depth": true})); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: figma-cli file get URL_OR_KEY [--depth N] [--raw]")
	}
	target, err := figmaurl.Parse(fs.Arg(0))
	if err != nil {
		return err
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	file, body, err := client.File(ctx, target.FileKey, *depth)
	if err != nil {
		return err
	}
	if *raw {
		return writeRaw(a.stdout, body)
	}
	printFileSummary(a.stdout, target.FileKey, file)
	return nil
}

func (a *app) runNode(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("node command required: inspect")
	}
	switch args[0] {
	case "inspect":
		return a.runNodeInspect(ctx, args[1:])
	default:
		return fmt.Errorf("unknown node command %q", args[0])
	}
}

func (a *app) runNodeInspect(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("node inspect", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	nodeID := fs.String("node", "", "Figma node ID; defaults to node-id from URL when present")
	raw := fs.Bool("raw", false, "print raw JSON")
	if err := fs.Parse(interspersed(args, map[string]bool{"node": true})); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: figma-cli node inspect URL_OR_KEY --node NODE_ID [--raw]")
	}
	target, err := figmaurl.Parse(fs.Arg(0))
	if err != nil {
		return err
	}
	id := firstNonEmpty(*nodeID, target.NodeID)
	if id == "" {
		return errors.New("--node is required when URL has no node-id")
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	resp, body, err := client.Nodes(ctx, target.FileKey, []string{id})
	if err != nil {
		return err
	}
	if *raw {
		return writeRaw(a.stdout, body)
	}
	wrap, ok := resp.Nodes[id]
	if !ok {
		return fmt.Errorf("node %s not found in response", id)
	}
	printNode(a.stdout, wrap.Document, 0, 3)
	return nil
}

func (a *app) runImage(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("image command required: export")
	}
	switch args[0] {
	case "export":
		return a.runImageExport(ctx, args[1:])
	default:
		return fmt.Errorf("unknown image command %q", args[0])
	}
}

func (a *app) runImageExport(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("image export", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	nodeID := fs.String("node", "", "Figma node ID; defaults to node-id from URL when present")
	format := fs.String("format", "png", "export format: png, jpg, svg, pdf")
	scale := fs.Float64("scale", 1, "export scale")
	outDir := fs.String("out", "", "download image into directory instead of printing URL only")
	raw := fs.Bool("raw", false, "print raw JSON")
	if err := fs.Parse(interspersed(args, map[string]bool{"node": true, "format": true, "scale": true, "out": true})); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: figma-cli image export URL_OR_KEY --node NODE_ID [--format png] [--scale 2] [--out DIR]")
	}
	target, err := figmaurl.Parse(fs.Arg(0))
	if err != nil {
		return err
	}
	id := firstNonEmpty(*nodeID, target.NodeID)
	if id == "" {
		return errors.New("--node is required when URL has no node-id")
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	resp, body, err := client.Images(ctx, target.FileKey, []string{id}, *format, *scale)
	if err != nil {
		return err
	}
	if *raw {
		return writeRaw(a.stdout, body)
	}
	imageURL := resp.Images[id]
	if imageURL == "" {
		return fmt.Errorf("no image URL returned for node %s", id)
	}
	if *outDir == "" {
		fmt.Fprintf(a.stdout, "%s\n", imageURL)
		return nil
	}
	outPath := filepath.Join(*outDir, safeFileName(id)+"."+strings.TrimPrefix(*format, "."))
	if err := client.Download(ctx, imageURL, outPath); err != nil {
		return err
	}
	fmt.Fprintf(a.stdout, "Wrote %s\n", outPath)
	return nil
}

func (a *app) runComments(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("comments command required: list")
	}
	switch args[0] {
	case "list":
		return a.runCommentsList(ctx, args[1:])
	default:
		return fmt.Errorf("unknown comments command %q", args[0])
	}
}

func (a *app) runCommentsList(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("comments list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	raw := fs.Bool("raw", false, "print raw JSON")
	if err := fs.Parse(interspersed(args, nil)); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return errors.New("usage: figma-cli comments list URL_OR_KEY [--raw]")
	}
	target, err := figmaurl.Parse(fs.Arg(0))
	if err != nil {
		return err
	}
	client, err := a.client()
	if err != nil {
		return err
	}
	resp, body, err := client.Comments(ctx, target.FileKey)
	if err != nil {
		return err
	}
	if *raw {
		return writeRaw(a.stdout, body)
	}
	if len(resp.Comments) == 0 {
		fmt.Fprintln(a.stdout, "No comments")
		return nil
	}
	for _, c := range resp.Comments {
		fmt.Fprintf(a.stdout, "%s | %s | %s\n", c.ID, c.User.Handle, oneLine(c.Message, 120))
	}
	return nil
}

func (a *app) client() (*api.Client, error) {
	token, _, err := credential.Resolve()
	if err != nil {
		return nil, err
	}
	c := api.New(token)
	if baseURL := strings.TrimRight(os.Getenv("FIGMA_API_BASE_URL"), "/"); baseURL != "" {
		c.BaseURL = baseURL
	}
	c.Verbose = a.verbose
	c.VerboseOut = a.stderr
	return c, nil
}

func (a *app) printHelp() {
	fmt.Fprint(a.stdout, `figma-cli is an agent-friendly Figma REST API CLI.

Usage:
  figma-cli init --token-stdin
  figma-cli me [--raw]
  figma-cli file get URL_OR_KEY [--depth N] [--raw]
  figma-cli node inspect URL_OR_KEY --node NODE_ID [--raw]
  figma-cli image export URL_OR_KEY --node NODE_ID [--format png] [--scale 2] [--out DIR]
  figma-cli comments list URL_OR_KEY [--raw]

Authentication:
  FIGMA_TOKEN overrides the OS keyring for one invocation.
  init stores a token in the OS keyring and never echoes it.

Global flags:
  -v, --verbose   print request diagnostics without token values
`)
}

func printFileSummary(w io.Writer, key string, file *api.File) {
	fmt.Fprintf(w, "Key: %s\nName: %s\nVersion: %s\nLast modified: %s\n", key, file.Name, file.Version, file.LastModified)
	fmt.Fprintf(w, "Root: %s (%s)\n", file.Document.Name, file.Document.Type)
	if len(file.Document.Children) > 0 {
		fmt.Fprintln(w, "Top-level nodes:")
		for _, child := range file.Document.Children {
			fmt.Fprintf(w, "- %s | %s | %s\n", child.ID, child.Type, child.Name)
		}
	}
}

func printNode(w io.Writer, n api.Node, depth, maxDepth int) {
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(w, "%s- %s | %s | %s", indent, n.ID, n.Type, n.Name)
	if n.AbsoluteBoundingBox != nil {
		fmt.Fprintf(w, " | %.0fx%.0f", n.AbsoluteBoundingBox.Width, n.AbsoluteBoundingBox.Height)
	}
	if n.Characters != "" {
		fmt.Fprintf(w, " | text=%q", oneLine(n.Characters, 80))
	}
	fmt.Fprintln(w)
	if depth >= maxDepth && len(n.Children) > 0 {
		fmt.Fprintf(w, "%s  ... %d child nodes hidden\n", indent, len(n.Children))
		return
	}
	for _, child := range n.Children {
		printNode(w, child, depth+1, maxDepth)
	}
}

func writeRaw(w io.Writer, body []byte) error {
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		_, werr := w.Write(body)
		return werr
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func oneLine(s string, max int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= max {
		return s
	}
	if max < 4 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func safeFileName(id string) string {
	replacer := strings.NewReplacer(":", "_", "/", "_", "\\", "_", " ", "_")
	return replacer.Replace(id)
}

func interspersed(args []string, valueFlags map[string]bool) []string {
	if len(args) == 0 {
		return args
	}
	var flags []string
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positional = append(positional, arg)
			continue
		}
		flags = append(flags, arg)
		name := strings.TrimLeft(arg, "-")
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			name = name[:eq]
		}
		if valueFlags[name] && !strings.Contains(arg, "=") && i+1 < len(args) {
			i++
			flags = append(flags, args[i])
		}
	}
	return append(flags, positional...)
}
