// Environment management JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log("Environment management loaded");

    // Add event listener for the "Add New Environment File" card
    const addEnvCard = document.querySelector('.add-card');
    if (addEnvCard) {
        addEnvCard.addEventListener('click', function(e) {
            e.preventDefault();
            showAddEnvModal();
        });
    }

    // Add event listeners for the "Add Variable" button
    const addVarBtn = document.querySelector('.add-var-btn');
    if (addVarBtn) {
        addVarBtn.addEventListener('click', function() {
            const envName = this.getAttribute('data-env');
            showAddVarModal(envName);
        });
    }

    // Add event listeners for edit variable buttons
    const editVarBtns = document.querySelectorAll('.edit-var-btn');
    editVarBtns.forEach(function(btn) {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const key = this.getAttribute('data-key');
            const value = this.getAttribute('data-value');
            const envName = document.querySelector('.add-var-btn').getAttribute('data-env');
            showEditVarModal(envName, key, value);
        });
    });

    // Add event listeners for delete variable buttons
    const deleteVarBtns = document.querySelectorAll('.delete-var-btn');
    deleteVarBtns.forEach(function(btn) {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            const key = this.getAttribute('data-key');
            const envName = document.querySelector('.add-var-btn').getAttribute('data-env');
            showDeleteVarModal(envName, key);
        });
    });

    // Add event listeners for download buttons
    const downloadBtns = document.querySelectorAll('.download-btn');
    downloadBtns.forEach(function(btn) {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            const envName = this.getAttribute('data-env');
            downloadEnvFile(envName);
        });
    });

    // Add event listeners for delete environment buttons
    const deleteEnvBtns = document.querySelectorAll('.delete-env-btn');
    deleteEnvBtns.forEach(function(btn) {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            const envName = this.getAttribute('data-env');
            showDeleteEnvModal(envName);
        });
    });
});

// Function to show the "Add Environment File" modal
function showAddEnvModal() {
    // Create modal HTML
    const modalHTML = `
        <div class="modal-overlay">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>Add New Environment File</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="add-env-form">
                        <div class="form-group">
                            <label for="env-name">Environment Name:</label>
                            <input type="text" id="env-name" name="env-name" placeholder="e.g., production" required>
                        </div>
                        <div class="form-group">
                            <label for="copy-from">Copy from existing (optional):</label>
                            <select id="copy-from" name="copy-from">
                                <option value="">-- None --</option>
                                ${getEnvOptionsHTML()}
                            </select>
                        </div>
                        <div class="form-actions">
                            <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                            <button type="submit" class="primary-btn">Create</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `;

    // Add modal to the DOM
    document.body.insertAdjacentHTML('beforeend', modalHTML);

    // Add event listeners
    const modal = document.querySelector('.modal-overlay');
    const closeBtn = modal.querySelector('.modal-close');
    const cancelBtn = modal.querySelector('.cancel-btn');
    const form = modal.querySelector('#add-env-form');

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const envName = document.getElementById('env-name').value;
        const copyFrom = document.getElementById('copy-from').value;

        // Validate the form
        if (!envName) {
            showErrorToast('Please enter an environment name');
            return;
        }

        // Check if the environment name is valid (alphanumeric, underscores, hyphens)
        if (!/^[a-zA-Z0-9_-]+$/.test(envName)) {
            showErrorToast('Environment name can only contain letters, numbers, underscores, and hyphens');
            return;
        }

        createEnvFile(envName, copyFrom);
        modal.remove();
    });
}

// Function to show the "Add Variable" modal
function showAddVarModal(envName) {
    // Create modal HTML
    const modalHTML = `
        <div class="modal-overlay">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>Add Variable to ${envName}</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="add-var-form">
                        <div class="form-group">
                            <label for="var-key">Key:</label>
                            <input type="text" id="var-key" name="var-key" placeholder="e.g., API_URL" required pattern="^[A-Za-z0-9_]+$">
                            <div class="form-hint">Use only letters, numbers, and underscores (no spaces or special characters)</div>
                        </div>
                        <div class="form-group">
                            <label for="var-value">Value:</label>
                            <input type="text" id="var-value" name="var-value" placeholder="e.g., https://api.example.com" required>
                        </div>
                        <div id="key-error" class="error-message" style="display: none;"></div>
                        <div class="form-actions">
                            <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                            <button type="submit" class="primary-btn">Add</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `;

    // Add modal to the DOM
    document.body.insertAdjacentHTML('beforeend', modalHTML);

    // Add event listeners
    const modal = document.querySelector('.modal-overlay');
    const closeBtn = modal.querySelector('.modal-close');
    const cancelBtn = modal.querySelector('.cancel-btn');
    const form = modal.querySelector('#add-var-form');
    const keyInput = document.getElementById('var-key');
    const keyError = document.getElementById('key-error');

    // Add input validation
    keyInput.addEventListener('input', function() {
        validateKey(keyInput, keyError);
    });

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();

        const key = keyInput.value;
        const value = document.getElementById('var-value').value;

        // Validate key format
        if (!validateKey(keyInput, keyError)) {
            return;
        }

        addVariable(envName, key, value);
        modal.remove();
    });
}

// Function to show the "Edit Variable" modal
function showEditVarModal(envName, key, value) {
    // Create modal HTML
    const modalHTML = `
        <div class="modal-overlay">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>Edit Variable in ${envName}</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <form id="edit-var-form">
                        <div class="form-group">
                            <label for="var-key">Key:</label>
                            <input type="text" id="var-key" name="var-key" value="${key}" readonly>
                            <div class="form-hint">Keys cannot be edited. Delete this variable and create a new one if needed.</div>
                        </div>
                        <div class="form-group">
                            <label for="var-value">Value:</label>
                            <input type="text" id="var-value" name="var-value" value="${value}" required>
                        </div>
                        <div class="form-actions">
                            <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                            <button type="submit" class="primary-btn">Update</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    `;

    // Add modal to the DOM
    document.body.insertAdjacentHTML('beforeend', modalHTML);

    // Add event listeners
    const modal = document.querySelector('.modal-overlay');
    const closeBtn = modal.querySelector('.modal-close');
    const cancelBtn = modal.querySelector('.cancel-btn');
    const form = modal.querySelector('#edit-var-form');

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const newValue = document.getElementById('var-value').value;

        updateVariable(envName, key, newValue);
        modal.remove();
    });
}

// Function to show the "Delete Variable" confirmation toast
function showDeleteVarModal(envName, key) {
    showConfirmationToast(
        `Are you sure you want to delete the variable "${key}" from the ${envName} environment?`,
        'Confirm Deletion',
        {
            confirmText: 'Delete',
            cancelText: 'Cancel',
            confirmButtonClass: 'primary-btn',
            onConfirm: () => {
                // Show loading toast
                const loadingToastId = showInfoToast('Deleting variable...', 'Please wait', 0);

                // Delete the variable
                deleteVariable(envName, key, loadingToastId);
            }
        }
    );
}

// Function to show the "Delete Environment" confirmation toast
function showDeleteEnvModal(envName) {
    showConfirmationToast(
        `Are you sure you want to delete the ${envName} environment file? This will permanently delete all variables in this environment.`,
        'Confirm Environment Deletion',
        {
            confirmText: 'Delete',
            cancelText: 'Cancel',
            confirmButtonClass: 'primary-btn',
            onConfirm: () => {
                // Show loading toast
                const loadingToastId = showInfoToast('Deleting environment...', 'Please wait', 0);

                // Delete the environment file
                deleteEnvFile(envName, loadingToastId);
            }
        }
    );
}

// Helper function to get HTML options for environment select
function getEnvOptionsHTML() {
    let options = '';
    const envCards = document.querySelectorAll('.info-card:not(.add-card):not(.empty-card)');

    envCards.forEach(function(card) {
        const envName = card.querySelector('.card-label').textContent;
        options += `<option value="${envName}">${envName}</option>`;
    });

    return options;
}

// These functions would be implemented to interact with the server
// For now, they just reload the page to show the changes
function createEnvFile(envName, copyFrom) {
    console.log(`Creating environment file: ${envName}, copy from: ${copyFrom}`);

    // Show loading toast
    const loadingToastId = showInfoToast('Creating environment file...', 'Please wait', 0);

    // Create a form to submit the request
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/api/environment/create';
    form.style.display = 'none';

    // Add the environment name
    const envNameInput = document.createElement('input');
    envNameInput.type = 'hidden';
    envNameInput.name = 'env_name';
    envNameInput.value = envName;
    form.appendChild(envNameInput);

    // Add the copy from parameter if provided
    if (copyFrom) {
        const copyFromInput = document.createElement('input');
        copyFromInput.type = 'hidden';
        copyFromInput.name = 'copy_from';
        copyFromInput.value = copyFrom;
        form.appendChild(copyFromInput);
    }

    // Add a callback to show success message
    const iframe = document.createElement('iframe');
    iframe.name = 'create-env-frame';
    iframe.style.display = 'none';
    document.body.appendChild(iframe);

    iframe.onload = function() {
        // Remove loading toast
        removeToast(loadingToastId);

        // Show success toast
        const message = copyFrom
            ? `Environment "${envName}" created successfully (copied from ${copyFrom})`
            : `Environment "${envName}" created successfully`;

        showSuccessToast(message);

        // Reload the page after a short delay
        setTimeout(() => {
            window.location.reload();
        }, 1500);
    };

    form.target = 'create-env-frame';

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function addVariable(envName, key, value) {
    console.log(`Adding variable to ${envName}: ${key}=${value}`);

    // Show loading toast
    const loadingToastId = showInfoToast('Adding variable...', 'Please wait', 0);

    // Create a form to submit the request
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/api/environment/add-variable';
    form.style.display = 'none';

    // Add the environment name
    const envNameInput = document.createElement('input');
    envNameInput.type = 'hidden';
    envNameInput.name = 'env_name';
    envNameInput.value = envName;
    form.appendChild(envNameInput);

    // Add the key
    const keyInput = document.createElement('input');
    keyInput.type = 'hidden';
    keyInput.name = 'key';
    keyInput.value = key;
    form.appendChild(keyInput);

    // Add the value
    const valueInput = document.createElement('input');
    valueInput.type = 'hidden';
    valueInput.name = 'value';
    valueInput.value = value;
    form.appendChild(valueInput);

    // Add a callback to show success message
    const iframe = document.createElement('iframe');
    iframe.name = 'add-var-frame';
    iframe.style.display = 'none';
    document.body.appendChild(iframe);

    iframe.onload = function() {
        // Remove loading toast
        removeToast(loadingToastId);

        // Show success toast
        showSuccessToast(`Variable "${key}" added successfully to ${envName} environment`);

        // Reload the page after a short delay
        setTimeout(() => {
            window.location.reload();
        }, 1500);
    };

    form.target = 'add-var-frame';

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function updateVariable(envName, key, value) {
    // For now, updating a variable is the same as adding it (upsert)
    addVariable(envName, key, value);
}

function deleteVariable(envName, key, loadingToastId) {
    console.log(`Deleting variable from ${envName}: ${key}`);

    // Create a form to submit the request
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/api/environment/delete-variable';
    form.style.display = 'none';

    // Add the environment name
    const envNameInput = document.createElement('input');
    envNameInput.type = 'hidden';
    envNameInput.name = 'env_name';
    envNameInput.value = envName;
    form.appendChild(envNameInput);

    // Add the key
    const keyInput = document.createElement('input');
    keyInput.type = 'hidden';
    keyInput.name = 'key';
    keyInput.value = key;
    form.appendChild(keyInput);

    // Add a callback to show success message
    const iframe = document.createElement('iframe');
    iframe.name = 'delete-var-frame';
    iframe.style.display = 'none';
    document.body.appendChild(iframe);

    iframe.onload = function() {
        // Remove loading toast if it exists
        if (loadingToastId) {
            removeToast(loadingToastId);
        }

        // Show success toast
        showSuccessToast(`Variable "${key}" deleted successfully from ${envName} environment`);

        // Reload the page after a short delay
        setTimeout(() => {
            window.location.reload();
        }, 1500);
    };

    form.target = 'delete-var-frame';

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function deleteEnvFile(envName, loadingToastId) {
    console.log(`Deleting environment file: ${envName}`);

    // Create a form to submit the request
    const form = document.createElement('form');
    form.method = 'POST';
    form.action = '/api/environment/delete-env';
    form.style.display = 'none';

    // Add the environment name
    const envNameInput = document.createElement('input');
    envNameInput.type = 'hidden';
    envNameInput.name = 'env_name';
    envNameInput.value = envName;
    form.appendChild(envNameInput);

    // Add a callback to show success message
    const iframe = document.createElement('iframe');
    iframe.name = 'delete-env-frame';
    iframe.style.display = 'none';
    document.body.appendChild(iframe);

    iframe.onload = function() {
        // Remove loading toast if it exists
        if (loadingToastId) {
            removeToast(loadingToastId);
        }

        // Show success toast
        showSuccessToast(`Environment "${envName}" deleted successfully`);

        // Reload the page after a short delay
        setTimeout(() => {
            window.location.reload();
        }, 1500);
    };

    form.target = 'delete-env-frame';

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function downloadEnvFile(envName) {
    console.log(`Downloading environment file: ${envName}`);

    // Create a link to download the file
    window.location.href = `/api/environment/download?env_name=${encodeURIComponent(envName)}`;
}

// Validate environment variable key format
function validateKey(inputElement, errorElement) {
    const key = inputElement.value;
    const keyRegex = /^[A-Za-z_][A-Za-z0-9_]*$/;

    if (!keyRegex.test(key)) {
        let errorMessage;
        if (/^\d/.test(key)) {
            errorMessage = "Key must not start with a number (Dart variable naming convention)";
        } else {
            errorMessage = "Key must contain only letters, numbers, and underscores (no spaces or special characters)";
        }

        // Show in the form error element
        if (errorElement) {
            errorElement.textContent = errorMessage;
            errorElement.style.display = "block";
        }

        // Also show a toast notification
        showWarningToast(errorMessage, "Invalid Key Format");

        if (inputElement) {
            inputElement.classList.add("input-error");
        }
        return false;
    } else {
        if (errorElement) {
            errorElement.style.display = "none";
        }
        if (inputElement) {
            inputElement.classList.remove("input-error");
        }
        return true;
    }
}
