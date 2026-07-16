#!/usr/bin/env python3
import http.server
import json
import os
import pathlib
import shutil
import socketserver
import subprocess
import tempfile
import threading


ROOT = pathlib.Path(__file__).resolve().parents[1]
BIN = ROOT / "bin" / "figma-cli"


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path.startswith("/assets/1_2.png"):
            data = b"\x89PNG\r\n\x1a\nmock"
            self.send_response(200)
            self.send_header("content-type", "image/png")
            self.send_header("content-length", str(len(data)))
            self.end_headers()
            self.wfile.write(data)
            return

        if self.headers.get("X-Figma-Token") != "mock-token":
            self.send_response(401)
            self.end_headers()
            return

        if self.path == "/v1/me":
            self.json({"id": "u1", "email": "designer@example.com", "handle": "Designer"})
        elif self.path.startswith("/v1/files/FILE123/nodes"):
            self.json({
                "nodes": {
                    "1:2": {
                        "document": {
                            "id": "1:2",
                            "name": "Hero",
                            "type": "FRAME",
                            "absoluteBoundingBox": {"x": 0, "y": 0, "width": 1440, "height": 900},
                            "children": [
                                {"id": "1:3", "name": "Title", "type": "TEXT", "characters": "Hello Figma"}
                            ],
                        }
                    }
                }
            })
        elif self.path.startswith("/v1/files/FILE123/comments"):
            self.json({"comments": [{"id": "c1", "message": "Looks good", "created_at": "2026-07-16T00:00:00Z", "user": {"id": "u1", "handle": "Designer"}}]})
        elif self.path.startswith("/v1/images/FILE123"):
            base = f"http://127.0.0.1:{self.server.server_address[1]}"
            self.json({"images": {"1:2": base + "/assets/1_2.png"}})
        elif self.path.startswith("/v1/files/FILE123/variables/local"):
            self.json({
                "meta": {
                    "variables": {
                        "v1": {"id": "v1", "name": "Color/Primary", "resolvedType": "COLOR", "variableCollectionId": "vc1", "scopes": ["ALL_SCOPES"]}
                    },
                    "variableCollections": {"vc1": {"id": "vc1", "name": "Colors"}},
                }
            })
        elif self.path.startswith("/v1/files/FILE123"):
            self.json({
                "name": "Mock File",
                "version": "1",
                "lastModified": "2026-07-16T00:00:00Z",
                "document": {"id": "0:0", "name": "Document", "type": "DOCUMENT", "children": [{"id": "1:2", "name": "Hero", "type": "FRAME"}]},
                "components": {"comp1": {"key": "comp1", "name": "Button/Primary", "description": "Primary button", "remote": False}},
                "styles": {"style1": {"key": "style1", "name": "Text/H1", "description": "Hero heading", "styleType": "TEXT", "remote": False}},
            })
        else:
            self.send_response(404)
            self.end_headers()

    def json(self, body):
        data = json.dumps(body).encode()
        self.send_response(200)
        self.send_header("content-type", "application/json")
        self.send_header("content-length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)

    def log_message(self, *args):
        pass


def run(cmd, env):
    out = subprocess.check_output(cmd, cwd=ROOT, env=env, text=True, stderr=subprocess.STDOUT)
    print("$", " ".join(cmd))
    print(out.strip())
    print()
    return out


def main():
    if not BIN.exists():
        subprocess.check_call(["make", "build"], cwd=ROOT)
    out_dir = pathlib.Path(tempfile.mkdtemp(prefix="figma-cli-mock-"))
    with socketserver.TCPServer(("127.0.0.1", 0), Handler) as srv:
        thread = threading.Thread(target=srv.serve_forever, daemon=True)
        thread.start()
        env = os.environ.copy()
        env["FIGMA_TOKEN"] = "mock-token"
        env["FIGMA_API_BASE_URL"] = f"http://127.0.0.1:{srv.server_address[1]}/v1"
        env["XDG_CACHE_HOME"] = str(out_dir / "cache")

        checks = [
            (["./bin/figma-cli", "me"], "Designer"),
            (["./bin/figma-cli", "file", "get", "FILE123", "--no-cache"], "Mock File"),
            (["./bin/figma-cli", "node", "inspect", "FILE123", "--node", "1:2", "--no-cache"], "Hello Figma"),
            (["./bin/figma-cli", "comments", "list", "FILE123", "--no-cache"], "Looks good"),
            (["./bin/figma-cli", "components", "list", "FILE123", "--no-cache"], "Button/Primary"),
            (["./bin/figma-cli", "styles", "list", "FILE123", "--no-cache"], "Text/H1"),
            (["./bin/figma-cli", "variables", "list", "FILE123", "--no-cache"], "Color/Primary"),
            (["./bin/figma-cli", "image", "export", "FILE123", "--node", "1:2", "--out", str(out_dir), "--no-cache"], "Wrote"),
        ]
        for cmd, expected in checks:
            output = run(cmd, env)
            if expected not in output:
                raise SystemExit(f"expected {expected!r} in output for {' '.join(cmd)}")
        exported = out_dir / "1_2.png"
        if not exported.exists():
            raise SystemExit(f"missing exported image {exported}")
        srv.shutdown()
    shutil.rmtree(out_dir, ignore_errors=True)
    print("mock verification passed")


if __name__ == "__main__":
    main()
