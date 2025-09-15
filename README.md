# GitCLI
A lightweight **command-line tool in Go** that generates an ASCII contribution heatmap of your Git commits, similar to GitHub’s contribution graph.

---

## ✨ Features
- Scan local Git repositories for commit history
- Filter commits by email address
- Visualize contributions in an ASCII heatmap (weeks × days)
- Easily add new repo folders to scan

---

## 🚀 Usage

### Add a folder to scan
```bash
go run main.go scan.go stats.go --add /path/to/repo --email "you@example.com"
