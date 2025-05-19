// JavaScript for Assets Page
document.addEventListener('DOMContentLoaded', function() {
    console.log("Assets page loaded");

    // Asset counts
    let assetCounts = {
        images: 0,
        animations: 0,
        audio: 0,
        videos: 0,
        json: 0,
        svgs: 0,
        misc: 0,
        total: 0
    };

    // Toggle summary section
    const toggleSummaryBtn = document.getElementById('toggle-summary');
    const summaryContent = document.getElementById('summary-content');

    if (toggleSummaryBtn && summaryContent) {
        toggleSummaryBtn.addEventListener('click', function() {
            summaryContent.classList.toggle('collapsed');
            toggleSummaryBtn.classList.toggle('collapsed');
        });
    }

    // Upload button toggle
    const uploadBtn = document.getElementById('upload-assets-btn');
    const uploadArea = document.getElementById('assets-upload-area');

    if (uploadBtn && uploadArea) {
        uploadBtn.addEventListener('click', function() {
            uploadArea.style.display = uploadArea.style.display === 'none' ? 'block' : 'none';
        });
    }

    // Migrate button functionality
    const migrateBtn = document.getElementById('migrate-assets-btn');
    const migrateDialog = document.getElementById('migrate-dialog');
    const migrateClose = document.getElementById('migrate-close');
    const migrateCancel = document.getElementById('migrate-cancel');
    const migrateConfirm = document.getElementById('migrate-confirm');

    if (migrateBtn && migrateDialog) {
        migrateBtn.addEventListener('click', function() {
            migrateDialog.style.display = 'flex';
        });

        if (migrateClose) {
            migrateClose.addEventListener('click', function() {
                migrateDialog.style.display = 'none';
            });
        }

        if (migrateCancel) {
            migrateCancel.addEventListener('click', function() {
                migrateDialog.style.display = 'none';
            });
        }

        if (migrateConfirm) {
            migrateConfirm.addEventListener('click', function() {
                migrateAssets();
            });
        }
    }

    // Drag and drop functionality
    const dropZone = document.getElementById('drop-zone');
    const fileInput = document.getElementById('asset-upload');
    const selectedFilesContainer = document.getElementById('selected-files');
    const uploadSubmitBtn = document.getElementById('upload-submit-btn');
    const uploadForm = document.getElementById('upload-form');

    let selectedFiles = [];

    if (dropZone && fileInput) {
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, preventDefaults, false);
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        ['dragenter', 'dragover'].forEach(eventName => {
            dropZone.addEventListener(eventName, highlight, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, unhighlight, false);
        });

        function highlight() {
            dropZone.classList.add('highlight');
        }

        function unhighlight() {
            dropZone.classList.remove('highlight');
        }

        dropZone.addEventListener('drop', handleDrop, false);

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            handleFiles(files);
        }

        fileInput.addEventListener('change', function() {
            handleFiles(this.files);
        });

        function handleFiles(files) {
            if (files.length === 0) return;

            // Add files to selected files array
            Array.from(files).forEach(file => {
                if (!selectedFiles.some(f => f.name === file.name && f.size === file.size)) {
                    selectedFiles.push(file);
                }
            });

            // Update selected files display
            updateSelectedFilesDisplay();

            // Show upload button
            if (uploadSubmitBtn) {
                uploadSubmitBtn.style.display = 'inline-flex';
            }
        }

        function updateSelectedFilesDisplay() {
            if (!selectedFilesContainer) return;

            selectedFilesContainer.innerHTML = '';

            selectedFiles.forEach((file, index) => {
                const fileElement = document.createElement('div');
                fileElement.className = 'selected-file';

                const fileName = document.createElement('span');
                fileName.className = 'file-name';
                fileName.textContent = file.name;

                const removeBtn = document.createElement('button');
                removeBtn.className = 'remove-file';
                removeBtn.innerHTML = '<i class="fas fa-times"></i>';
                removeBtn.addEventListener('click', function() {
                    selectedFiles.splice(index, 1);
                    updateSelectedFilesDisplay();

                    if (selectedFiles.length === 0 && uploadSubmitBtn) {
                        uploadSubmitBtn.style.display = 'none';
                    }
                });

                fileElement.appendChild(fileName);
                fileElement.appendChild(removeBtn);
                selectedFilesContainer.appendChild(fileElement);
            });
        }

        // Handle form submission
        if (uploadForm) {
            uploadForm.addEventListener('submit', function(e) {
                e.preventDefault();
                uploadFiles();
            });
        }
    }

    // Load assets on page load
    loadAssets();

    // Search functionality
    const searchInput = document.getElementById('asset-search');

    if (searchInput) {
        searchInput.addEventListener('input', function() {
            const searchTerm = this.value.toLowerCase();
            const assetRows = document.querySelectorAll('.asset-row');

            assetRows.forEach(row => {
                const assetName = row.querySelector('.asset-name').textContent.toLowerCase();
                const assetPath = row.querySelector('.asset-path').textContent.toLowerCase();

                if (assetName.includes(searchTerm) || assetPath.includes(searchTerm)) {
                    row.style.display = '';
                } else {
                    row.style.display = 'none';
                }
            });
        });
    }

    // Filter functionality
    const filterSelect = document.getElementById('asset-type-filter');

    if (filterSelect) {
        filterSelect.addEventListener('change', function() {
            const filterValue = this.value;
            const assetRows = document.querySelectorAll('.asset-row');

            if (filterValue === 'all') {
                assetRows.forEach(row => {
                    row.style.display = '';
                });
            } else {
                assetRows.forEach(row => {
                    const assetType = row.querySelector('.asset-type').textContent.toLowerCase();

                    if (assetType === filterValue ||
                        (filterValue === 'images' && assetType === 'image')) {
                        row.style.display = '';
                    } else {
                        row.style.display = 'none';
                    }
                });
            }
        });
    }

    // Function to load assets from the server
    function loadAssets() {
        const tableBody = document.getElementById('assets-table-body');
        const loadingIndicator = document.getElementById('loading-assets');
        const noAssetsMessage = document.getElementById('no-assets-message');

        if (!tableBody || !loadingIndicator || !noAssetsMessage) return;

        // Show loading indicator
        loadingIndicator.style.display = 'block';
        noAssetsMessage.style.display = 'none';
        tableBody.innerHTML = '';

        // Reset asset counts
        assetCounts = {
            images: 0,
            animations: 0,
            audio: 0,
            videos: 0,
            json: 0,
            svgs: 0,
            misc: 0,
            total: 0
        };

        // Fetch assets from the server
        fetch('/api/assets/list')
            .then(response => response.json())
            .then(data => {
                // Hide loading indicator
                loadingIndicator.style.display = 'none';

                if (data.success && data.assets) {
                    let totalAssets = 0;

                    // Process each asset type
                    for (const [type, files] of Object.entries(data.assets)) {
                        // Update asset counts
                        assetCounts[type] = files.length;
                        totalAssets += files.length;

                        // Add rows for each asset
                        files.forEach(file => {
                            const row = createAssetRow(file, type);
                            tableBody.appendChild(row);
                        });
                    }

                    // Update total count
                    assetCounts.total = totalAssets;

                    // Update count displays
                    updateAssetCounts();

                    // Show no assets message if no assets found
                    if (totalAssets === 0) {
                        noAssetsMessage.style.display = 'block';
                    }
                } else {
                    console.error('Failed to load assets:', data);
                    noAssetsMessage.style.display = 'block';
                }
            })
            .catch(error => {
                console.error('Error loading assets:', error);
                loadingIndicator.style.display = 'none';
                noAssetsMessage.style.display = 'block';
            });
    }

    // Function to create an asset row
    function createAssetRow(fileName, assetType) {
        const row = document.createElement('tr');
        row.className = 'asset-row';

        // Preview cell
        const previewCell = document.createElement('td');
        previewCell.className = 'asset-preview';

        let iconClass = 'fa-file';

        // Set icon based on asset type
        switch (assetType) {
            case 'images':
                iconClass = 'fa-file-image';
                break;
            case 'animations':
                iconClass = 'fa-film';
                break;
            case 'audio':
                iconClass = 'fa-file-audio';
                break;
            case 'videos':
                iconClass = 'fa-file-video';
                break;
            case 'json':
                iconClass = 'fa-file-code';
                break;
            case 'svgs':
                iconClass = 'fa-bezier-curve';
                break;
            case 'misc':
                iconClass = 'fa-file-alt';
                break;
        }

        previewCell.innerHTML = `<i class="fas ${iconClass}"></i>`;

        // Name cell
        const nameCell = document.createElement('td');
        nameCell.className = 'asset-name';
        nameCell.textContent = fileName;

        // Type cell
        const typeCell = document.createElement('td');
        typeCell.className = 'asset-type';
        typeCell.textContent = assetType;

        // Path cell
        const pathCell = document.createElement('td');
        pathCell.className = 'asset-path';
        pathCell.textContent = `assets/${assetType}/${fileName}`;

        // Actions cell
        const actionsCell = document.createElement('td');
        actionsCell.className = 'asset-actions';

        const downloadBtn = document.createElement('button');
        downloadBtn.className = 'table-btn';
        downloadBtn.title = 'Download';
        downloadBtn.innerHTML = '<i class="fas fa-download"></i>';
        downloadBtn.addEventListener('click', function() {
            downloadAsset(fileName, assetType);
        });

        const deleteBtn = document.createElement('button');
        deleteBtn.className = 'table-btn delete-btn';
        deleteBtn.title = 'Delete';
        deleteBtn.innerHTML = '<i class="fas fa-trash"></i>';
        deleteBtn.addEventListener('click', function() {
            deleteAsset(fileName, assetType, row);
        });

        actionsCell.appendChild(downloadBtn);
        actionsCell.appendChild(deleteBtn);

        // Add cells to row
        row.appendChild(previewCell);
        row.appendChild(nameCell);
        row.appendChild(typeCell);
        row.appendChild(pathCell);
        row.appendChild(actionsCell);

        return row;
    }

    // Function to update asset counts in the UI
    function updateAssetCounts() {
        document.getElementById('total-assets-count').textContent = assetCounts.total;
        document.getElementById('images-count').textContent = assetCounts.images;
        document.getElementById('animations-count').textContent = assetCounts.animations;
        document.getElementById('audio-count').textContent = assetCounts.audio;
        document.getElementById('videos-count').textContent = assetCounts.videos;
        document.getElementById('json-count').textContent = assetCounts.json;
        document.getElementById('svgs-count').textContent = assetCounts.svgs;
        document.getElementById('misc-count').textContent = assetCounts.misc;
    }

    // Function to upload files
    function uploadFiles() {
        if (selectedFiles.length === 0) return;

        const formData = new FormData();
        const assetType = document.getElementById('asset-type-select').value;

        // Add asset type if specified
        if (assetType) {
            formData.append('asset_type', assetType);
        }

        // Add files to form data
        selectedFiles.forEach(file => {
            formData.append('files', file);
        });

        // Show loading state
        const uploadSubmitBtn = document.getElementById('upload-submit-btn');
        if (uploadSubmitBtn) {
            uploadSubmitBtn.disabled = true;
            uploadSubmitBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Uploading...';
        }

        // Upload files to server
        fetch('/api/assets/upload', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Reset selected files
                selectedFiles = [];
                updateSelectedFilesDisplay();

                // Hide upload button
                if (uploadSubmitBtn) {
                    uploadSubmitBtn.style.display = 'none';
                    uploadSubmitBtn.disabled = false;
                    uploadSubmitBtn.innerHTML = '<i class="fas fa-upload"></i> Upload Selected Files';
                }

                // Hide upload area
                if (uploadArea) {
                    uploadArea.style.display = 'none';
                }

                // Reload assets
                loadAssets();

                // Show success message
                alert('Files uploaded successfully!');
            } else {
                console.error('Failed to upload files:', data);
                alert('Failed to upload files. Please try again.');

                // Reset upload button
                if (uploadSubmitBtn) {
                    uploadSubmitBtn.disabled = false;
                    uploadSubmitBtn.innerHTML = '<i class="fas fa-upload"></i> Upload Selected Files';
                }
            }
        })
        .catch(error => {
            console.error('Error uploading files:', error);
            alert('Error uploading files. Please try again.');

            // Reset upload button
            if (uploadSubmitBtn) {
                uploadSubmitBtn.disabled = false;
                uploadSubmitBtn.innerHTML = '<i class="fas fa-upload"></i> Upload Selected Files';
            }
        });
    }

    // Function to download an asset
    function downloadAsset(fileName, assetType) {
        window.location.href = `/api/assets/download?asset_name=${encodeURIComponent(fileName)}&asset_type=${encodeURIComponent(assetType)}`;
    }

    // Function to delete an asset
    function deleteAsset(fileName, assetType, row) {
        if (!confirm(`Are you sure you want to delete the asset "${fileName}"?`)) {
            return;
        }

        // Log the values being sent to help debug
        console.log('Deleting asset:', { fileName, assetType });

        const formData = new FormData();
        formData.append('asset_name', fileName);
        formData.append('asset_type', assetType);

        fetch('/api/assets/delete', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Remove row from table
                if (row) {
                    row.remove();
                }

                // Update asset counts
                assetCounts[assetType]--;
                assetCounts.total--;
                updateAssetCounts();

                // Show no assets message if no assets left
                if (assetCounts.total === 0) {
                    document.getElementById('no-assets-message').style.display = 'block';
                }

                // Show success message
                alert('Asset deleted successfully!');
            } else {
                console.error('Failed to delete asset:', data);
                alert('Failed to delete asset. Please try again.');
            }
        })
        .catch(error => {
            console.error('Error deleting asset:', error);
            alert('Error deleting asset. Please try again.');
        });
    }

    // Function to migrate assets
    function migrateAssets() {
        // Hide migrate dialog
        const migrateDialog = document.getElementById('migrate-dialog');
        if (migrateDialog) {
            migrateDialog.style.display = 'none';
        }

        // Show loading state
        const migrateBtn = document.getElementById('migrate-assets-btn');
        if (migrateBtn) {
            migrateBtn.disabled = true;
            migrateBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Migrating...';
        }

        // Migrate assets on server
        fetch('/api/assets/migrate', {
            method: 'POST'
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Reload assets
                loadAssets();

                // Show success message
                alert('Assets migrated successfully!');
            } else {
                console.error('Failed to migrate assets:', data);
                alert('Failed to migrate assets. Please try again.');
            }

            // Reset migrate button
            if (migrateBtn) {
                migrateBtn.disabled = false;
                migrateBtn.innerHTML = '<i class="fas fa-exchange-alt"></i> Migrate Assets';
            }
        })
        .catch(error => {
            console.error('Error migrating assets:', error);
            alert('Error migrating assets. Please try again.');

            // Reset migrate button
            if (migrateBtn) {
                migrateBtn.disabled = false;
                migrateBtn.innerHTML = '<i class="fas fa-exchange-alt"></i> Migrate Assets';
            }
        });
    }
});
