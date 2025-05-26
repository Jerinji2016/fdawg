// Namer functionality
class NamerManager {
    constructor() {
        this.platforms = [];
        this.currentNames = {};
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadPlatforms();
        this.loadCurrentNames();
    }

    bindEvents() {
        // Refresh button
        document.getElementById('refresh-names-btn').addEventListener('click', () => {
            this.loadCurrentNames();
            this.loadPlatforms();
        });

        // Toggle platform information section
        document.getElementById('toggle-platform-info').addEventListener('click', () => {
            this.togglePlatformInfoSection();
        });

        // Universal name setting
        document.getElementById('set-universal-btn').addEventListener('click', () => {
            this.setUniversalName();
        });
    }

    async loadPlatforms() {
        try {
            const response = await fetch('/api/namer/platforms');
            if (!response.ok) {
                throw new Error('Failed to load platforms');
            }

            const data = await response.json();
            this.platforms = data.all || [];
            this.renderPlatformForms();
        } catch (error) {
            console.error('Error loading platforms:', error);
            showToast('Failed to load platforms', 'error');
        }
    }

    async loadCurrentNames() {
        this.showLoading();

        try {
            const response = await fetch('/api/namer/get', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({})
            });

            if (!response.ok) {
                throw new Error('Failed to load current names');
            }

            const data = await response.json();
            this.currentNames = data;
            // Re-render platform forms with updated current names
            this.renderPlatformForms();
        } catch (error) {
            console.error('Error loading current names:', error);
            showToast('Failed to load current app names', 'error');
        } finally {
            this.hideLoading();
        }
    }



    renderPlatformForms() {
        const container = document.getElementById('platform-forms');

        const html = this.platforms.map(platform => {
            const platformIcon = this.getPlatformIcon(platform.id);
            const isAvailable = platform.available;

            // Find current app name for this platform
            const currentAppName = this.currentNames.app_names ?
                this.currentNames.app_names.find(app => app.platform === platform.id) : null;

            const currentDisplayName = currentAppName ? currentAppName.display_name || 'Not set' : 'Loading...';
            const hasInternalName = currentAppName && currentAppName.internal_name &&
                currentAppName.internal_name !== currentAppName.display_name;

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
                        <div class="current-name-display">
                            <div class="current-name-header">
                                <div class="current-name-label">Current Name:</div>
                                <button class="edit-toggle-btn" onclick="namerManager.toggleEdit('${platform.id}')" title="Edit app name">
                                    <i class="fas fa-edit"></i>
                                </button>
                            </div>
                            <div class="current-name-value">${currentDisplayName}</div>
                            ${hasInternalName ? `<div class="current-internal-name">(Internal: ${currentAppName.internal_name})</div>` : ''}
                            ${currentAppName && currentAppName.error ? `<div class="error-text"><i class="fas fa-exclamation-triangle"></i> ${currentAppName.error}</div>` : ''}
                        </div>

                        <div class="edit-name-section" id="edit-section-${platform.id}" style="display: none;">
                            <input
                                type="text"
                                id="platform-${platform.id}"
                                class="platform-input"
                                placeholder="Enter new app name for ${platform.name}"
                                value=""
                            >
                            <div class="edit-buttons">
                                <button class="save-platform-btn" onclick="namerManager.savePlatformName('${platform.id}')">
                                    <i class="fas fa-save"></i> Save
                                </button>
                                <button class="cancel-edit-btn" onclick="namerManager.cancelEdit('${platform.id}')">
                                    <i class="fas fa-times"></i>
                                </button>
                            </div>
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

    async setUniversalName() {
        const universalName = document.getElementById('universal-name').value.trim();

        if (!universalName) {
            showToast('Please enter an app name', 'warning');
            return;
        }

        const availablePlatforms = this.platforms.filter(p => p.available);
        if (availablePlatforms.length === 0) {
            showToast('No available platforms found', 'warning');
            return;
        }

        const message = `Set "${universalName}" as the app name for all available platforms?`;
        const details = availablePlatforms.map(p => `• ${p.name}`).join('<br>');

        this.showConfirmationDialog(message, details, async () => {
            await this.executeSetNames({ universal: universalName });
        });
    }

    toggleEdit(platformId) {
        const editSection = document.getElementById(`edit-section-${platformId}`);
        const isVisible = editSection.style.display !== 'none';

        if (isVisible) {
            editSection.style.display = 'none';
        } else {
            editSection.style.display = 'block';
            // Focus on the input field
            const input = document.getElementById(`platform-${platformId}`);
            if (input) {
                input.focus();
            }
        }
    }

    cancelEdit(platformId) {
        const editSection = document.getElementById(`edit-section-${platformId}`);
        const input = document.getElementById(`platform-${platformId}`);

        // Hide edit section and clear input
        editSection.style.display = 'none';
        if (input) {
            input.value = '';
        }
    }

    async savePlatformName(platformId) {
        const input = document.getElementById(`platform-${platformId}`);
        const newName = input.value.trim();

        if (!newName) {
            showToast('Please enter an app name', 'warning');
            return;
        }

        const platform = this.platforms.find(p => p.id === platformId);
        const platformName = platform ? platform.name : this.capitalizeFirst(platformId);

        const message = `Set "${newName}" as the app name for ${platformName}?`;
        const details = `Platform: ${platformName}<br>New Name: "${newName}"`;

        this.showConfirmationDialog(message, details, async () => {
            await this.executeSetNames({ platforms: { [platformId]: newName } });
            // Hide edit section and clear input after successful save
            this.cancelEdit(platformId);
        });
    }

    async executeSetNames(request) {
        this.showLoading();

        try {
            const response = await fetch('/api/namer/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to set app names');
            }

            await response.json();
            showToast('App names updated successfully!', 'success');

            // Refresh the current names display
            await this.loadCurrentNames();

            // Clear forms
            this.clearAllForms();

        } catch (error) {
            console.error('Error setting app names:', error);
            showToast(`Failed to set app names: ${error.message}`, 'error');
        } finally {
            this.hideLoading();
        }
    }

    clearAllForms() {
        document.getElementById('universal-name').value = '';
        // Clear all platform inputs
        this.platforms.forEach(platform => {
            const input = document.getElementById(`platform-${platform.id}`);
            if (input) {
                input.value = '';
            }
        });
        showToast('All forms cleared', 'info');
    }

    togglePlatformInfoSection() {
        const container = document.getElementById('platform-info-container');
        const toggleBtn = document.getElementById('toggle-platform-info');
        const icon = toggleBtn.querySelector('i');

        if (container.style.display === 'none') {
            container.style.display = 'grid';
            icon.classList.remove('fa-chevron-down');
            icon.classList.add('fa-chevron-up');
        } else {
            container.style.display = 'none';
            icon.classList.remove('fa-chevron-up');
            icon.classList.add('fa-chevron-down');
        }
    }

    showConfirmationDialog(message, details, action) {
        const modalHTML = `
            <div class="modal-overlay">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3><i class="fas fa-question-circle"></i> Confirm Action</h3>
                        <button class="modal-close">&times;</button>
                    </div>
                    <div class="modal-body">
                        <p>${message}</p>
                        <div class="confirmation-details">${details}</div>
                    </div>
                    <div class="modal-footer">
                        <button class="primary-btn confirm-yes-btn">
                            <i class="fas fa-check"></i> Yes, Continue
                        </button>
                        <button class="secondary-btn confirm-no-btn">
                            <i class="fas fa-times"></i> Cancel
                        </button>
                    </div>
                </div>
            </div>
        `;

        document.body.insertAdjacentHTML('beforeend', modalHTML);
        const modal = document.querySelector('.modal-overlay:last-child');

        // Add event listeners
        modal.querySelector('.modal-close').addEventListener('click', () => this.hideConfirmationDialog(modal));
        modal.querySelector('.confirm-no-btn').addEventListener('click', () => this.hideConfirmationDialog(modal));
        modal.querySelector('.confirm-yes-btn').addEventListener('click', async () => {
            this.hideConfirmationDialog(modal);
            await action();
        });

        // Close on overlay click
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                this.hideConfirmationDialog(modal);
            }
        });
    }

    hideConfirmationDialog(modal) {
        if (modal) {
            modal.remove();
        }
    }

    showLoading() {
        if (!document.getElementById('loading-overlay')) {
            const loadingHTML = `
                <div id="loading-overlay" class="modal-overlay">
                    <div class="loading-spinner">
                        <i class="fas fa-spinner fa-spin"></i>
                        <p>Processing...</p>
                    </div>
                </div>
            `;
            document.body.insertAdjacentHTML('beforeend', loadingHTML);
        } else {
            document.getElementById('loading-overlay').style.display = 'flex';
        }
    }

    hideLoading() {
        const loadingOverlay = document.getElementById('loading-overlay');
        if (loadingOverlay) {
            loadingOverlay.style.display = 'none';
        }
    }

    getPlatformIcon(platform) {
        const icons = {
            android: 'fab fa-android',
            ios: 'fab fa-apple',
            macos: 'fab fa-apple',
            linux: 'fab fa-linux',
            windows: 'fab fa-windows',
            web: 'fas fa-globe'
        };
        return icons[platform] || 'fas fa-desktop';
    }

    capitalizeFirst(str) {
        return str.charAt(0).toUpperCase() + str.slice(1);
    }
}

// Global instance for onclick handlers
let namerManager;

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    namerManager = new NamerManager();
});
