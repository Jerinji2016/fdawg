// Bundler functionality
class BundlerManager {
    constructor() {
        this.platforms = [];
        this.currentBundleIDs = {};
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadPlatforms();
        this.loadCurrentBundleIDs();
    }

    bindEvents() {
        // Refresh button
        document.getElementById('refresh-bundleids-btn').addEventListener('click', () => {
            this.loadCurrentBundleIDs();
            this.loadPlatforms();
        });

        // Toggle platform information section
        document.getElementById('toggle-platform-info').addEventListener('click', () => {
            this.toggleSection('platform-info-container', 'toggle-platform-info');
        });

        // Toggle guidelines section
        document.getElementById('toggle-guidelines').addEventListener('click', () => {
            this.toggleSection('guidelines-container', 'toggle-guidelines');
        });

        // Universal bundle ID setting
        document.getElementById('set-universal-btn').addEventListener('click', () => {
            this.setUniversalBundleID();
        });

        // Real-time validation on input
        document.getElementById('universal-bundleid').addEventListener('input', () => {
            this.clearValidationResult();
        });
    }

    toggleSection(containerId, toggleBtnId) {
        const container = document.getElementById(containerId);
        const toggleBtn = document.getElementById(toggleBtnId);
        const icon = toggleBtn.querySelector('i');

        if (container.style.display === 'none') {
            container.style.display = 'block';
            icon.className = 'fas fa-chevron-up';
        } else {
            container.style.display = 'none';
            icon.className = 'fas fa-chevron-down';
        }
    }

    async loadPlatforms() {
        try {
            const response = await fetch('/api/bundler/platforms');
            if (!response.ok) {
                throw new Error('Failed to load platforms');
            }

            const data = await response.json();
            this.platforms = data.all || [];
            this.renderPlatformForms();
        } catch (error) {
            console.error('Error loading platforms:', error);
            showToast('Failed to load platform information', 'error');
        }
    }

    async loadCurrentBundleIDs() {
        this.showLoading();

        try {
            const response = await fetch('/api/bundler/get', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({})
            });

            if (!response.ok) {
                throw new Error('Failed to load current bundle IDs');
            }

            const data = await response.json();
            this.currentBundleIDs = data;
            // Re-render platform forms with updated current bundle IDs
            this.renderPlatformForms();
        } catch (error) {
            console.error('Error loading current bundle IDs:', error);
            showToast('Failed to load current bundle IDs', 'error');
        } finally {
            this.hideLoading();
        }
    }

    renderPlatformForms() {
        const container = document.getElementById('platform-forms');

        // Platform icon mapping
        const platformIcons = {
            'android': 'fab fa-android',
            'ios': 'fab fa-apple',
            'macos': 'fab fa-apple',
            'linux': 'fab fa-linux',
            'windows': 'fab fa-windows',
            'web': 'fas fa-globe'
        };

        const html = this.platforms.map(platform => {
            const isAvailable = platform.available;
            const currentBundleID = this.currentBundleIDs.bundle_ids?.find(b => b.platform === platform.id);
            const currentID = currentBundleID?.bundle_id || '';
            const currentNamespace = currentBundleID?.namespace || '';
            const hasNamespace = currentNamespace && currentNamespace !== currentID;
            const platformIcon = platformIcons[platform.id] || 'fas fa-desktop';

            return `
                <div class="platform-form-card ${isAvailable ? '' : 'unavailable'}">
                    <div class="platform-header">
                        <div class="platform-title">
                            <i class="${platformIcon}"></i>
                            ${platform.name}
                        </div>
                        <span class="platform-status ${isAvailable ? 'available' : 'unavailable'}">
                            ${isAvailable ? 'Available' : 'Not Available'}
                        </span>
                    </div>

                    ${isAvailable ? `
                        <div class="current-bundleid-display">
                            <div class="current-bundleid-header">
                                <div class="current-bundleid-label">Current Bundle ID:</div>
                                <button class="edit-toggle-btn" onclick="bundlerManager.toggleEdit('${platform.id}')" title="Edit bundle ID">
                                    <i class="fas fa-edit"></i>
                                </button>
                            </div>
                            <div class="current-bundleid-value">${currentID}</div>
                            ${hasNamespace ? `<div class="current-namespace">(Namespace: ${currentNamespace})</div>` : ''}
                            ${currentBundleID && currentBundleID.error ? `<div class="error-text"><i class="fas fa-exclamation-triangle"></i> ${currentBundleID.error}</div>` : ''}
                        </div>

                        <div class="edit-bundleid-section" id="edit-section-${platform.id}" style="display: none;">
                            <input
                                type="text"
                                id="platform-${platform.id}"
                                class="platform-input"
                                placeholder="Enter new bundle ID for ${platform.name}"
                                value=""
                            >
                            <div class="edit-buttons">
                                <button class="save-platform-btn" onclick="bundlerManager.savePlatformBundleID('${platform.id}')">
                                    <i class="fas fa-save"></i> Save
                                </button>
                                <button class="cancel-edit-btn" onclick="bundlerManager.cancelEdit('${platform.id}')">
                                    <i class="fas fa-times"></i>
                                </button>
                            </div>
                            <div id="platform-validation-${platform.id}" class="validation-result" style="display: none;"></div>
                        </div>
                    ` : `
                        <div class="unavailable-message">
                            <i class="fas fa-info-circle"></i> Platform not available in this project
                        </div>
                    `}
                </div>
            `;
        }).join('');

        container.innerHTML = html;
    }

    async validateBundleIDAPI(bundleID) {
        try {
            const response = await fetch('/api/bundler/validate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ bundle_id: bundleID })
            });

            if (!response.ok) {
                throw new Error('Failed to validate bundle ID');
            }

            const data = await response.json();
            return {
                valid: data.valid,
                error: data.error || 'Bundle ID is valid'
            };
        } catch (error) {
            console.error('Error validating bundle ID:', error);
            return {
                valid: false,
                error: 'Failed to validate bundle ID'
            };
        }
    }

    showValidationResult(elementId, isValid, message) {
        const element = document.getElementById(elementId);
        element.style.display = 'block';
        element.className = `validation-result ${isValid ? 'valid' : 'invalid'}`;
        element.innerHTML = `
            <i class="fas ${isValid ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i>
            ${message}
        `;
    }

    clearValidationResult() {
        const element = document.getElementById('universal-validation-result');
        element.style.display = 'none';
    }

    async setUniversalBundleID() {
        const universalBundleID = document.getElementById('universal-bundleid').value.trim();

        if (!universalBundleID) {
            showToast('Please enter a bundle ID', 'warning');
            return;
        }

        // Validate bundle ID format before proceeding
        const validation = await this.validateBundleIDAPI(universalBundleID);
        if (!validation.valid) {
            this.showValidationResult('universal-validation-result', false, validation.error);
            showToast('Invalid bundle ID format', 'error');
            return;
        }

        // Clear any previous validation messages
        this.clearValidationResult();

        const availablePlatforms = this.platforms.filter(p => p.available);
        if (availablePlatforms.length === 0) {
            showToast('No available platforms found', 'warning');
            return;
        }

        const message = `Set "${universalBundleID}" as the bundle ID for all available platforms?`;
        const details = availablePlatforms.map(p => `• ${p.name}`).join('<br>');

        this.showConfirmationDialog(message, details, async () => {
            await this.executeSetBundleIDs({ universal: universalBundleID });
        });
    }

    toggleEdit(platformId) {
        const editSection = document.getElementById(`edit-section-${platformId}`);
        const isVisible = editSection.style.display !== 'none';

        if (isVisible) {
            this.cancelEdit(platformId);
        } else {
            editSection.style.display = 'block';
            const input = document.getElementById(`platform-${platformId}`);
            input.focus();
        }
    }

    cancelEdit(platformId) {
        const editSection = document.getElementById(`edit-section-${platformId}`);
        const input = document.getElementById(`platform-${platformId}`);
        const validationResult = document.getElementById(`platform-validation-${platformId}`);
        
        editSection.style.display = 'none';
        input.value = '';
        if (validationResult) {
            validationResult.style.display = 'none';
        }
    }

    async savePlatformBundleID(platformId) {
        const input = document.getElementById(`platform-${platformId}`);
        const newBundleID = input.value.trim();

        if (!newBundleID) {
            showToast('Please enter a bundle ID', 'warning');
            return;
        }

        // Validate bundle ID format before proceeding
        const validation = await this.validateBundleIDAPI(newBundleID);
        if (!validation.valid) {
            this.showValidationResult(`platform-validation-${platformId}`, false, validation.error);
            showToast('Invalid bundle ID format', 'error');
            return;
        }

        // Clear any previous validation messages
        const validationElement = document.getElementById(`platform-validation-${platformId}`);
        if (validationElement) {
            validationElement.style.display = 'none';
        }

        const platform = this.platforms.find(p => p.id === platformId);
        const platformName = platform ? platform.name : this.capitalizeFirst(platformId);

        const message = `Set "${newBundleID}" as the bundle ID for ${platformName}?`;
        const details = `Platform: ${platformName}<br>New Bundle ID: "${newBundleID}"`;

        this.showConfirmationDialog(message, details, async () => {
            await this.executeSetBundleIDs({ platforms: { [platformId]: newBundleID } });
            // Hide edit section and clear input after successful save
            this.cancelEdit(platformId);
        });
    }

    async executeSetBundleIDs(request) {
        this.showLoading();

        try {
            const response = await fetch('/api/bundler/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to set bundle IDs');
            }

            await response.json();
            showToast('Bundle IDs updated successfully!', 'success');

            // Refresh the current bundle IDs display
            await this.loadCurrentBundleIDs();

            // Clear forms
            this.clearAllForms();

        } catch (error) {
            console.error('Error setting bundle IDs:', error);
            showToast(`Failed to set bundle IDs: ${error.message}`, 'error');
        } finally {
            this.hideLoading();
        }
    }

    showConfirmationDialog(message, details, onConfirm) {
        const dialog = document.createElement('div');
        dialog.className = 'confirmation-dialog-overlay';
        dialog.innerHTML = `
            <div class="confirmation-dialog">
                <div class="dialog-header">
                    <h4><i class="fas fa-question-circle"></i> Confirm Bundle ID Change</h4>
                </div>
                <div class="dialog-content">
                    <p>${message}</p>
                    <div class="dialog-details">${details}</div>
                </div>
                <div class="dialog-actions">
                    <button class="cancel-btn" onclick="this.closest('.confirmation-dialog-overlay').remove()">
                        <i class="fas fa-times"></i> Cancel
                    </button>
                    <button class="confirm-btn" onclick="bundlerManager.confirmAction(this, arguments[0])">
                        <i class="fas fa-check"></i> Confirm
                    </button>
                </div>
            </div>
        `;

        // Store the callback function
        dialog.querySelector('.confirm-btn').onConfirmCallback = onConfirm;

        document.body.appendChild(dialog);
    }

    async confirmAction(button) {
        const dialog = button.closest('.confirmation-dialog-overlay');
        const onConfirm = button.onConfirmCallback;

        dialog.remove();

        if (onConfirm) {
            await onConfirm();
        }
    }

    clearAllForms() {
        // Clear universal input
        document.getElementById('universal-bundleid').value = '';
        this.clearValidationResult();

        // Clear all platform inputs
        this.platforms.forEach(platform => {
            const input = document.getElementById(`platform-${platform.id}`);
            if (input) {
                input.value = '';
            }
            this.cancelEdit(platform.id);
        });
    }

    capitalizeFirst(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
    }

    showLoading() {
        // Show loading indicator
        const loadingIndicator = document.createElement('div');
        loadingIndicator.id = 'bundler-loading';
        loadingIndicator.className = 'loading-overlay';
        loadingIndicator.innerHTML = `
            <div class="loading-spinner">
                <i class="fas fa-spinner fa-spin"></i>
                <span>Processing bundle IDs...</span>
            </div>
        `;
        document.body.appendChild(loadingIndicator);
    }

    hideLoading() {
        // Hide loading indicator
        const loadingIndicator = document.getElementById('bundler-loading');
        if (loadingIndicator) {
            loadingIndicator.remove();
        }
    }
}

// Initialize bundler manager when DOM is loaded
let bundlerManager;
document.addEventListener('DOMContentLoaded', function() {
    bundlerManager = new BundlerManager();
});
