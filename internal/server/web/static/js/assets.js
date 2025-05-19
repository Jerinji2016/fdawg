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

    // Browse files button
    const browseBtn = document.getElementById('browse-assets-btn');
    const fileInput = document.getElementById('asset-upload');

    if (browseBtn && fileInput) {
        browseBtn.addEventListener('click', function() {
            fileInput.click();
        });
    }

    // Migrate button functionality
    const migrateBtn = document.getElementById('migrate-assets-btn');

    if (migrateBtn) {
        migrateBtn.addEventListener('click', function() {
            showConfirmationToast(
                'This will organize your assets into folders by type (images, animations, audio, etc.) and update your pubspec.yaml file.',
                'Migrate Assets?',
                {
                    confirmText: 'Migrate',
                    cancelText: 'Cancel',
                    confirmButtonClass: 'primary-btn',
                    onConfirm: () => {
                        migrateAssets();
                    }
                }
            );
        });
    }

    // Page-level drag and drop functionality
    const pageDropZone = document.getElementById('page-drop-zone');
    const floatingContainer = document.getElementById('floating-files-container');
    const floatingFilesList = document.getElementById('floating-files-list');
    const uploadSelectedBtn = document.getElementById('upload-selected-files');
    const closeFloatingBtn = document.getElementById('close-floating-container');
    const assetsContainer = document.getElementById('assets-management-container');

    let selectedFiles = [];

    // Prevent defaults for drag events
    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    // Setup page-level drag and drop
    if (assetsContainer && pageDropZone && floatingContainer) {
        // Add event listeners to the entire document
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            document.addEventListener(eventName, preventDefaults, false);
        });

        // Show drop zone when files are dragged over the document
        ['dragenter', 'dragover'].forEach(eventName => {
            document.addEventListener(eventName, function(e) {
                // Only show if it's a file being dragged
                if (e.dataTransfer.types.includes('Files')) {
                    pageDropZone.classList.add('active');
                }
            }, false);
        });

        // Hide drop zone when files are dragged out or dropped
        document.addEventListener('dragleave', function(e) {
            // Only hide if dragleave is on the document itself and not a child element
            if (e.target === document.documentElement) {
                pageDropZone.classList.remove('active');
            }
        }, false);

        // Handle file drop
        document.addEventListener('drop', function(e) {
            pageDropZone.classList.remove('active');

            // Only process if files were dropped
            if (e.dataTransfer.files.length > 0) {
                handleFiles(e.dataTransfer.files);
            }
        }, false);

        // Handle file input change
        if (fileInput) {
            fileInput.addEventListener('change', function() {
                if (this.files.length > 0) {
                    handleFiles(this.files);
                }
            });
        }

        // Handle files
        function handleFiles(files) {
            if (files.length === 0) return;

            // Add files to selected files array
            Array.from(files).forEach(file => {
                if (!selectedFiles.some(f => f.name === file.name && f.size === file.size)) {
                    selectedFiles.push(file);
                }
            });

            // Update floating files display
            updateFloatingFilesDisplay();

            // Show floating container
            floatingContainer.classList.add('active');
        }

        // Update floating files display
        function updateFloatingFilesDisplay() {
            if (!floatingFilesList) return;

            floatingFilesList.innerHTML = '';

            if (selectedFiles.length === 0) {
                const emptyMessage = document.createElement('div');
                emptyMessage.className = 'empty-message';
                emptyMessage.textContent = 'No files selected';
                floatingFilesList.appendChild(emptyMessage);
                return;
            }

            selectedFiles.forEach((file, index) => {
                const fileElement = document.createElement('div');
                fileElement.className = 'floating-file';

                const fileInfo = document.createElement('div');
                fileInfo.className = 'file-info';

                // Set icon based on file type
                let iconClass = 'fa-file';
                const ext = file.name.split('.').pop().toLowerCase();

                if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'bmp'].includes(ext)) {
                    iconClass = 'fa-file-image';
                } else if (['mp3', 'wav', 'ogg'].includes(ext)) {
                    iconClass = 'fa-file-audio';
                } else if (['mp4', 'webm', 'avi', 'mov'].includes(ext)) {
                    iconClass = 'fa-file-video';
                } else if (ext === 'json') {
                    iconClass = 'fa-file-code';
                } else if (ext === 'svg') {
                    iconClass = 'fa-bezier-curve';
                }

                const icon = document.createElement('i');
                icon.className = `fas ${iconClass}`;
                fileInfo.appendChild(icon);

                const fileName = document.createElement('span');
                fileName.className = 'file-name';
                fileName.textContent = file.name;
                fileInfo.appendChild(fileName);

                const removeBtn = document.createElement('button');
                removeBtn.className = 'remove-file';
                removeBtn.innerHTML = '<i class="fas fa-times"></i>';
                removeBtn.addEventListener('click', function() {
                    selectedFiles.splice(index, 1);
                    updateFloatingFilesDisplay();

                    if (selectedFiles.length === 0) {
                        floatingContainer.classList.remove('active');
                    }
                });

                fileElement.appendChild(fileInfo);
                fileElement.appendChild(removeBtn);
                floatingFilesList.appendChild(fileElement);
            });
        }

        // Close floating container
        if (closeFloatingBtn) {
            closeFloatingBtn.addEventListener('click', function() {
                floatingContainer.classList.remove('active');
                selectedFiles = [];
            });
        }

        // Upload selected files
        if (uploadSelectedBtn) {
            uploadSelectedBtn.addEventListener('click', function() {
                if (selectedFiles.length > 0) {
                    uploadFiles();
                }
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
        const assetType = document.getElementById('floating-asset-type').value;

        // Add asset type if specified
        if (assetType) {
            formData.append('asset_type', assetType);
        }

        // Add files to form data
        selectedFiles.forEach(file => {
            formData.append('files', file);
        });

        // Show loading state
        if (uploadSelectedBtn) {
            uploadSelectedBtn.disabled = true;
            uploadSelectedBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Uploading...';
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

                // Hide floating container
                if (floatingContainer) {
                    floatingContainer.classList.remove('active');
                }

                // Reload assets
                loadAssets();

                // Show success message
                showSuccessToast('Files uploaded successfully!');
            } else {
                console.error('Failed to upload files:', data);
                showErrorToast('Failed to upload files. Please try again.');

                // Reset upload button
                if (uploadSelectedBtn) {
                    uploadSelectedBtn.disabled = false;
                    uploadSelectedBtn.innerHTML = '<i class="fas fa-upload"></i> Upload';
                }
            }
        })
        .catch(error => {
            console.error('Error uploading files:', error);
            showErrorToast('Error uploading files. Please try again.');

            // Reset upload button
            if (uploadSelectedBtn) {
                uploadSelectedBtn.disabled = false;
                uploadSelectedBtn.innerHTML = '<i class="fas fa-upload"></i> Upload';
            }
        });
    }

    // Function to download an asset
    function downloadAsset(fileName, assetType) {
        window.location.href = `/api/assets/download?asset_name=${encodeURIComponent(fileName)}&asset_type=${encodeURIComponent(assetType)}`;
    }

    // Function to delete an asset
    function deleteAsset(fileName, assetType, row) {
        showConfirmationToast(
            `Are you sure you want to delete "${fileName}"?`,
            'Confirm Deletion',
            {
                confirmText: 'Delete',
                cancelText: 'Cancel',
                confirmButtonClass: 'primary-btn',
                onConfirm: () => {
                    performAssetDeletion(fileName, assetType, row);
                }
            }
        );
    }

    // Function to perform the actual asset deletion
    function performAssetDeletion(fileName, assetType, row) {
        // Log the values being sent to help debug
        console.log('Deleting asset:', { fileName, assetType });

        const formData = new FormData();
        formData.append('asset_name', fileName);
        formData.append('asset_type', assetType);

        // Show a loading toast
        const loadingToastId = showInfoToast('Deleting asset...', 'Please wait', 0);

        fetch('/api/assets/delete', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            // Remove loading toast
            removeToast(loadingToastId);

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
                showSuccessToast('Asset deleted successfully!');
            } else {
                console.error('Failed to delete asset:', data);
                showErrorToast('Failed to delete asset. Please try again.');
            }
        })
        .catch(error => {
            // Remove loading toast
            removeToast(loadingToastId);

            console.error('Error deleting asset:', error);
            showErrorToast('Error deleting asset. Please try again.');
        });
    }

    // Function to migrate assets
    function migrateAssets() {
        // Show loading state
        const migrateBtn = document.getElementById('migrate-assets-btn');
        if (migrateBtn) {
            migrateBtn.disabled = true;
            migrateBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Migrating...';
        }

        // Show a loading toast
        const loadingToastId = showInfoToast('Migrating assets...', 'Please wait', 0);

        // Migrate assets on server
        fetch('/api/assets/migrate', {
            method: 'POST'
        })
        .then(response => response.json())
        .then(data => {
            // Remove loading toast
            removeToast(loadingToastId);

            if (data.success) {
                // Reload assets
                loadAssets();

                // Show success message
                showSuccessToast('Assets migrated successfully!');
            } else {
                console.error('Failed to migrate assets:', data);
                showErrorToast('Failed to migrate assets. Please try again.');
            }

            // Reset migrate button
            if (migrateBtn) {
                migrateBtn.disabled = false;
                migrateBtn.innerHTML = '<i class="fas fa-exchange-alt"></i> Migrate Assets';
            }
        })
        .catch(error => {
            // Remove loading toast
            removeToast(loadingToastId);

            console.error('Error migrating assets:', error);
            showErrorToast('Error migrating assets. Please try again.');

            // Reset migrate button
            if (migrateBtn) {
                migrateBtn.disabled = false;
                migrateBtn.innerHTML = '<i class="fas fa-exchange-alt"></i> Migrate Assets';
            }
        });
    }
});
