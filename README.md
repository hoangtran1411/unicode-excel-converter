# VNI/TCVN3 to Unicode Excel Converter

A high-performance Desktop Application built with **Go (Golang)** and **Wails**, designed to convert legacy Vietnamese encodings (VNI-Windows, TCVN3/ABC) in Excel files to standard Unicode.

<div align="center">
    <img src="https://img.shields.io/github/actions/workflow/status/hoangtran1411/convert-vni-to-unicode/ci.yml?style=flat-square&logo=github&label=Build" alt="CI Status" />
    <img src="https://img.shields.io/github/v/release/hoangtran1411/convert-vni-to-unicode?style=flat-square&logo=github" alt="Release" />
    <img src="https://img.shields.io/github/license/hoangtran1411/convert-vni-to-unicode?style=flat-square" alt="License" />
</div>

## 🚀 Features

- **Format Preservation**: Keeps your Excel styling intact!
    - Preserves **Bold**, *Italic*, Underline.
    - Preserves Font Sizes and Colors.
    - **Smart Font Mapping**: Automatically maps legacy fonts to Unicode equivalents (e.g., `.VnTime` -> `Times New Roman`, `VNI-Times` -> `Times New Roman`).
    - **Default Font**: Enforces `Arial` for converted text if no specific map is found.
- **Dual Encoding Support**:
    - **VNI-Windows**: Detects and converts headers and content using VNI fonts (e.g., `VNI-Times`).
    - **TCVN3 (ABC)**: Detects and converts TCVN3 fonts (e.g., `.VnTime`, `.VnArial`).
    - **Auto-Detection**: Smartly detects encoding based on Font Name and Content heuristics.
- **High Performance**:
    - Multi-threaded processing using a **Worker Pool** pattern.
    - Handles large Excel files without freezing the UI.
- **Modern UI**:
    - Premium Dark Theme with Glassmorphism effects.
    - Drag & Drop file support.
    - Real-time progress bar.
- **Auto-Update**:
    - Automatically checks for updates from GitHub Releases.
    - One-click in-app update.

## 🛠️ Technology Stack

<div align="center">
	<img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
	<img src="https://img.shields.io/badge/Wails-CC3534?style=for-the-badge&logo=wails&logoColor=white" alt="Wails" />
	<img src="https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white" alt="HTML5" />
	<img src="https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white" alt="CSS3" />
	<img src="https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black" alt="JavaScript" />
</div>

- **Backend**: Go (1.21+) - Robust and high-performance.
- **Frontend**: 
  - HTML5 / CSS3 (Glassmorphism UI)
  - Vanilla JavaScript (No heavy frameworks)
- **GUI Framework**: [Wails v2](https://wails.io) - Lightweight Application Runtime.
- **Excel Engine**: [Excelize v2](https://github.com/xuri/excelize) - High-performance Excel library.
- **CI/CD**: GitHub Actions & Dependabot.

## 📦 Installation

### From Source

1. **Prerequisites**:
    - Go 1.21 or later
    - Node.js (for building frontend assets)
    - Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

2. **Clone & Build**:
    ```bash
    git clone https://github.com/hoangtran1411/convert-vni-to-unicode.git
    cd convert-vni-to-unicode
    wails build
    ```

3. **Run**:
    The executable will be in `build/bin/convert-vni-to-unicode.exe`.

## 📖 Usage

1. Open the application.
2. **Drag & Drop** your Excel file (`.xlsx`) into the dotted area, or click "Browse File".
3. (Optional) Enter a specific Sheet Name. If left empty, it scans the first sheet (or all, depending on implementation).
4. Select **Source Encoding** (Auto-detect is recommended).
5. Click **START CONVERSION**.
6. The converted file will be saved in the **same folder** with the suffix `_output_yyyy_MM_dd_ss.xlsx`.

## 🧪 Development

### Running in Dev Mode
```bash
wails dev
```
This runs the app with hot-reload enabled for both Frontend and Backend.

### Running Tests
```bash
go test ./... -v
```

### Mocking Data for Test
To generate a sample Excel file with VNI/TCVN3 fonts for testing:
```bash
go run scripts/generate_sample.go
```
*(Note: You need to create this script or use the provided integration tests)*

## 🚀 Release Process

This project includes a **GitHub Actions** workflow to automate releases.

1. **Commit your changes**: Ensure `updater.go` has the correct `GitHubOwner` and `GitHubRepo`.
2. **Tag a new version**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. **Wait for Action**: GitHub will automatically build the Windows executable and create a new Release with the `.exe` file attached.
4. **Auto-Update**: Users running older versions will receive a notification to update to this new version.

## ⚙️ CI/CD

- **CI (`ci.yml`)**: Runs unit tests and linter on every Push and Pull Request to ensures code quality.
- **Release (`release.yml`)**: Builds the Windows binary and creates a GitHub Release when a new tag (e.g., `v1.0.0`) is pushed.
- **Dependabot**: Automatically checks for dependency updates weekly.

## 🏗️ Architecture

- **`main.go` / `app.go`**: Entry point and Wails binding boundaries.
- **`internal/engine`**:
    - `processor.go`: Core logic, manages Worker Pool and File I/O.
    - `format_preserver.go`: Handles formatting retention and font swapping.
    - `detector.go`: Heuristics for encoding detection.
- **`internal/converter`**: Pure Go logic for string conversion (VNI/TCVN3 maps).
- **`updater.go`**: Logic for self-update mechanism via GitHub API.

## 📝 License

MIT License.
