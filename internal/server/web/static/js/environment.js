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
                            <input type="text" id="var-key" name="var-key" placeholder="e.g., API_URL" required>
                        </div>
                        <div class="form-group">
                            <label for="var-value">Value:</label>
                            <input type="text" id="var-value" name="var-value" placeholder="e.g., https://api.example.com" required>
                        </div>
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

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    form.addEventListener('submit', function(e) {
        e.preventDefault();
        const key = document.getElementById('var-key').value;
        const value = document.getElementById('var-value').value;

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

// Function to show the "Delete Variable" confirmation modal
function showDeleteVarModal(envName, key) {
    // Create modal HTML
    const modalHTML = `
        <div class="modal-overlay">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>Delete Variable</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <p>Are you sure you want to delete the variable "${key}" from the ${envName} environment?</p>
                    <div class="form-actions">
                        <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                        <button type="button" class="primary-btn delete-btn">Delete</button>
                    </div>
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
    const deleteBtn = modal.querySelector('.delete-btn');

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    deleteBtn.addEventListener('click', function() {
        deleteVariable(envName, key);
        modal.remove();
    });
}

// Function to show the "Delete Environment" confirmation modal
function showDeleteEnvModal(envName) {
    // Create modal HTML
    const modalHTML = `
        <div class="modal-overlay">
            <div class="modal-content">
                <div class="modal-header">
                    <h3>Delete Environment</h3>
                    <button class="modal-close">&times;</button>
                </div>
                <div class="modal-body">
                    <p>Are you sure you want to delete the ${envName} environment file?</p>
                    <p class="warning-text">This will permanently delete all variables in this environment.</p>
                    <div class="form-actions">
                        <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                        <button type="button" class="primary-btn delete-btn">Delete</button>
                    </div>
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
    const deleteBtn = modal.querySelector('.delete-btn');

    closeBtn.addEventListener('click', function() {
        modal.remove();
    });

    cancelBtn.addEventListener('click', function() {
        modal.remove();
    });

    deleteBtn.addEventListener('click', function() {
        deleteEnvFile(envName);
        modal.remove();
    });
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

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function addVariable(envName, key, value) {
    console.log(`Adding variable to ${envName}: ${key}=${value}`);

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

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function updateVariable(envName, key, value) {
    // For now, updating a variable is the same as adding it (upsert)
    addVariable(envName, key, value);
}

function deleteVariable(envName, key) {
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

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function deleteEnvFile(envName) {
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

    // Submit the form
    document.body.appendChild(form);
    form.submit();
}

function downloadEnvFile(envName) {
    console.log(`Downloading environment file: ${envName}`);

    // Create a link to download the file
    window.location.href = `/api/environment/download?env_name=${encodeURIComponent(envName)}`;
}
