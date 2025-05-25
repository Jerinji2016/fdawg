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
        });

        // Toggle current names section
        document.getElementById('toggle-current-names').addEventListener('click', () => {
            this.toggleCurrentNamesSection();
        });

        // Universal name setting
        document.getElementById('set-universal-btn').addEventListener('click', () => {
            this.setUniversalName();
        });

        // Platform-specific name setting
        document.getElementById('set-platform-specific-btn').addEventListener('click', () => {
            this.setPlatformSpecificNames();
        });

        // Clear platform forms
        document.getElementById('clear-platform-forms-btn').addEventListener('click', () => {
            this.clearPlatformForms();
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
            this.renderCurrentNames();
        } catch (error) {
            console.error('Error loading current names:', error);
            showToast('Failed to load current app names', 'error');
        } finally {
            this.hideLoading();
        }
    }

    renderCurrentNames() {
        const container = document.getElementById('current-names-container');

        if (!this.currentNames.app_names || this.currentNames.app_names.length === 0) {
            container.innerHTML = '<tr><td colspan="3" class="empty-state"><i class="fas fa-info-circle"></i> No app names found</td></tr>';
            return;
        }

        const html = this.currentNames.app_names.map(appName => {
            const platformIcon = this.getPlatformIcon(appName.platform);
            const statusClass = appName.available ? 'available' : 'unavailable';
            const statusText = appName.available ? 'Available' : 'Not Available';

            let appNameDisplay = '';
            if (!appName.available) {
                appNameDisplay = '<span class="not-available-text">Not Available</span>';
            } else {
                const displayName = appName.display_name || 'Not set';
                const hasInternalName = appName.internal_name && appName.internal_name !== appName.display_name;

                appNameDisplay = `
                    <div class="app-name-cell">
                        <div class="name-display">${displayName}</div>
                        ${hasInternalName ? `<div class="name-internal">(Internal: ${appName.internal_name})</div>` : ''}
                        ${appName.error ? `<div class="error-text" title="${appName.error}"><i class="fas fa-exclamation-triangle"></i> ${appName.error}</div>` : ''}
                    </div>
                `;
            }

            return `
                <tr class="${appName.available ? '' : 'unavailable-row'}">
                    <td class="platform-cell">
                        <i class="${platformIcon}"></i>
                        <span class="platform-name">${this.capitalizeFirst(appName.platform)}</span>
                    </td>
                    <td class="app-name-cell">
                        ${appNameDisplay}
                    </td>
                    <td class="status-cell">
                        <span class="status-badge ${statusClass}">${statusText}</span>
                    </td>
                </tr>
            `;
        }).join('');

        container.innerHTML = html;
    }



    renderPlatformForms() {
        const container = document.getElementById('platform-forms');

        const html = this.platforms.map(platform => {
            const platformIcon = this.getPlatformIcon(platform.id);
            const isAvailable = platform.available;

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
                    <input
                        type="text"
                        id="platform-${platform.id}"
                        class="platform-input"
                        placeholder="Enter app name for ${platform.name}"
                        ${isAvailable ? '' : 'disabled'}
                    >
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

    async setPlatformSpecificNames() {
        const platformNames = {};
        let hasNames = false;

        this.platforms.forEach(platform => {
            if (platform.available) {
                const input = document.getElementById(`platform-${platform.id}`);
                const value = input.value.trim();
                if (value) {
                    platformNames[platform.id] = value;
                    hasNames = true;
                }
            }
        });

        if (!hasNames) {
            showToast('Please enter at least one app name', 'warning');
            return;
        }

        const message = 'Set the following platform-specific app names?';
        const details = Object.entries(platformNames)
            .map(([platform, name]) => `• ${this.capitalizeFirst(platform)}: "${name}"`)
            .join('<br>');

        this.showConfirmationDialog(message, details, async () => {
            await this.executeSetNames({ platforms: platformNames });
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

    clearPlatformForms() {
        this.platforms.forEach(platform => {
            const input = document.getElementById(`platform-${platform.id}`);
            if (input) {
                input.value = '';
            }
        });
        showToast('Platform forms cleared', 'info');
    }

    clearAllForms() {
        document.getElementById('universal-name').value = '';
        this.clearPlatformForms();
    }

    toggleCurrentNamesSection() {
        const summary = document.getElementById('current-names-summary');
        const toggleBtn = document.getElementById('toggle-current-names');
        const icon = toggleBtn.querySelector('i');

        if (summary.style.display === 'none') {
            summary.style.display = 'block';
            icon.classList.remove('fa-chevron-down');
            icon.classList.add('fa-chevron-up');
        } else {
            summary.style.display = 'none';
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

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new NamerManager();
});
