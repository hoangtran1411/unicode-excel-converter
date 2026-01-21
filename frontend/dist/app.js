// Main JS Logic

// Wails Runtime should be available via window.runtime or window.go
// If running in browser without Wails (dev mode), this mimics functionality.

const dropZone = document.getElementById('dropZone');
const fileInfo = document.getElementById('fileInfo');
const fileNameDisplay = document.getElementById('fileName');
const convertBtn = document.getElementById('convertBtn');
const progressContainer = document.getElementById('progressContainer');
const progressFill = document.getElementById('progressFill');
const progressText = document.getElementById('progressText');

let selectedPath = "";

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    // Check for updates
    checkForUpdates();
});

// Update Logic
let updateUrl = "";

async function checkForUpdates() {
    if (!window.go || !window.go.main) return;

    try {
        const info = await window.go.main.App.CheckForUpdate();
        if (info.available) {
            document.getElementById('new-version').textContent = info.latestVersion;
            document.getElementById('update-bar').style.display = 'flex';
            updateUrl = info.downloadUrl;
        }
    } catch (e) {
        console.error("Update check failed:", e);
    }
}

window.performUpdate = async () => {
    if (!updateUrl) return;

    const btn = document.querySelector('.btn-update');
    btn.textContent = "Downloading...";
    btn.disabled = true;

    try {
        await window.go.main.App.PerformUpdate(updateUrl);
    } catch (e) {
        showToast("Update failed: " + e, "error");
        btn.textContent = "Retry Update";
        btn.disabled = false;
    }
};

window.hideUpdate = () => {
    document.getElementById('update-bar').style.display = 'none';
};


// File Selection
function updateUIFileSelected(path) {
    if (path) {
        selectedPath = path;
        // Extract filename from path (simple split for display)
        const name = path.split(/[\\/]/).pop();
        fileNameDisplay.textContent = name;

        fileInfo.style.display = 'flex';
        convertBtn.disabled = false;

        // Hide "browse" button text somewhat? No stays same.
    } else {
        selectedPath = "";
        fileInfo.style.display = 'none';
        convertBtn.disabled = true;
    }
}

window.selectFile = async () => {
    try {
        // Call Go Backend
        const path = await window.go.main.App.SelectFile();
        updateUIFileSelected(path);
    } catch (e) {
        console.error(e);
    }
};

window.clearFile = (event) => {
    if (event) {
        event.stopPropagation();
    }
    updateUIFileSelected("");
};

// Start Conversion
window.startConversion = async () => {
    if (!selectedPath) return;

    try {
        convertBtn.disabled = true;
        convertBtn.textContent = "CONVERTING...";
        progressContainer.style.display = 'block';
        progressFill.style.width = '0%';
        progressText.textContent = "Initializing...";

        const sheetName = document.getElementById('sheetName').value;
        const encoding = document.getElementById('encoding').value;

        // Reset progress monitoring
        // We listen to "progress" event

        // Call Go
        const config = {
            inputPath: selectedPath,
            sheetName: sheetName,
            // Encoding is handled by Auto-detect usually, but if we passed it to config?
            // Updated App.go to include Encoding in Config?
            // Current Config struct in App.go: InputPath, SheetName.
            // Encoding logic is in Processor.DetectEncoding.
            // Ideally we should pass encoding preference if user forced it.
            // For now, let's assume Auto. If user selects VNI, we might want to force it.
            // TODO: Update backend to accept encoding override if needed.
        };

        const result = await window.go.main.App.Process(config);

        if (result.success) {
            progressFill.style.width = '100%';
            progressText.textContent = "Completed!";
            showToast(result.message, "success");
            // Optional: Show "Open Folder" button
        } else {
            progressFill.style.background = 'var(--danger)';
            progressText.textContent = "Failed";
            showToast("Error: " + result.message, "error");
        }
    } catch (e) {
        showToast("System Error: " + e, "error");
    } finally {
        convertBtn.disabled = false;
        convertBtn.textContent = "START CONVERSION";
    }
};

// Events from Backend
if (window.runtime) {
    window.runtime.EventsOn("progress", (count) => {
        // We don't know total, so we animate based on count or just show activity
        // Or if we implemented percentage in backend.
        // Current backend sends "processed count" (float64).
        // Let's just create a moving bar or display count.
        progressText.textContent = `Processed ${count} cells...`;

        // Fake visual progress if we don't know total:
        // Use an indeterminate animation or just fill slowly.
        // For now, let's just show the text and keep bar at 50% pulsing?
        // Or simpler: just text.
        progressFill.style.width = '50%';
    });

    window.runtime.EventsOn("updateProgress", (msg) => {
        showToast(msg, "info");
    });
}

// Toast Notification
function showToast(message, type = "info") {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.innerHTML = `<span>${message}</span>`;

    container.appendChild(toast);

    // Auto remove
    setTimeout(() => {
        toast.style.opacity = '0';
        setTimeout(() => toast.remove(), 300);
    }, 4000); // 4 seconds
}

// Drag & Drop
dropZone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropZone.classList.add('drag-over');
});

dropZone.addEventListener('dragleave', () => {
    dropZone.classList.remove('drag-over');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropZone.classList.remove('drag-over');

    const files = e.dataTransfer.files;
    if (files.length > 0) {
        const file = files[0];
        if (file.name.endsWith('.xlsx')) {
            // We need the full path. Browser security might block this in pure web,
            // but Wails WebView usually allows getting path if dropped?
            // Actually, Chrome/WebView DnD often gives File object but NOT full path.
            // Wails has specific DnD support or workaround.
            // Standard JS File API in WebView *might* not give `file.path`.
            // User might need to click Browse.
            // BUT: Wails 2 supports Drag and Drop on the Window triggering an event.
            // Or we try `file.path` (Electron style).
            // Let's try `file.path` if available (some webviews expose it).
            // If not, we might need a Wails Runtime 'OnDirDrop' handler.
            /* 
               Wails feature: runtime.OnFileDrop(func(x, y int, paths []string))
               We should implement that in Go (startup) and emit event to JS.
               But for now, let's rely on Browse button as primary.
            */
            // Attempt to read path (works in some Electron/Wails setups)
            // If 'path' property exists:
            // @ts-ignore
            if (file.path) {
                updateUIFileSelected(file.path);
            } else {
                showToast("Drag & Drop path detection not supported; please use Browse.", "error");
            }
        } else {
            showToast("Please select an .xlsx file", "error");
        }
    }
});
