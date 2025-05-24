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

        // Confirmation dialog
        document.getElementById('confirm-yes-btn').addEventListener('click', () => {
            this.executeConfirmedAction();
        });

        document.getElementById('confirm-no-btn').addEventListener('click', () => {
            this.hideConfirmationDialog();
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
            container.innerHTML = '<p class="text-muted">No app names found</p>';
            return;
        }

        const html = this.currentNames.app_names.map(appName => {
            const platformIcon = this.getPlatformIcon(appName.platform);
            const statusClass = appName.available ? 'available' : 'unavailable';
            const statusText = appName.available ? 'Available' : 'Not Available';

            return `
                <div class="platform-card ${appName.available ? '' : 'unavailable'}">
                    <div class="platform-header">
                        <i class="${platformIcon}"></i>
                        <h3>${this.capitalizeFirst(appName.platform)}</h3>
                        <span class="platform-status ${statusClass}">${statusText}</span>
                    </div>
                    ${appName.available ? this.renderPlatformNames(appName) : ''}
                    ${appName.error ? `<p class="text-error">${appName.error}</p>` : ''}
                </div>
            `;
        }).join('');

        container.innerHTML = html;
    }

    renderPlatformNames(appName) {
        let html = '<div class="platform-names">';
        
        if (appName.display_name) {
            html += `
                <div class="name-item">
                    <div class="name-label">Display Name:</div>
                    <div class="name-value">${appName.display_name}</div>
                </div>
            `;
        }

        if (appName.internal_name && appName.internal_name !== appName.display_name) {
            html += `
                <div class="name-item">
                    <div class="name-label">Internal Name:</div>
                    <div class="name-value">${appName.internal_name}</div>
                </div>
            `;
        }

        html += '</div>';
        return html;
    }

    renderPlatformForms() {
        const container = document.getElementById('platform-forms');
        
        const html = this.platforms.map(platform => {
            const platformIcon = this.getPlatformIcon(platform.id);
            const isAvailable = platform.available;

            return `
                <div class="platform-form ${isAvailable ? '' : 'unavailable'}">
                    <div class="platform-header">
                        <i class="${platformIcon}"></i>
                        <h4>${platform.name}</h4>
                        <span class="platform-status ${isAvailable ? 'available' : 'unavailable'}">
                            ${isAvailable ? 'Available' : 'Not Available'}
                        </span>
                    </div>
                    <input 
                        type="text" 
                        id="platform-${platform.id}" 
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
        const details = availablePlatforms.map(p => `• ${p.name}`).join('\n');
        
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
            .join('\n');

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

            const result = await response.json();
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

    showConfirmationDialog(message, details, action) {
        document.getElementById('confirmation-message').textContent = message;
        document.getElementById('confirmation-details').textContent = details;
        document.getElementById('confirmation-dialog').style.display = 'flex';
        this.pendingAction = action;
    }

    hideConfirmationDialog() {
        document.getElementById('confirmation-dialog').style.display = 'none';
        this.pendingAction = null;
    }

    async executeConfirmedAction() {
        this.hideConfirmationDialog();
        if (this.pendingAction) {
            await this.pendingAction();
        }
    }

    showLoading() {
        document.getElementById('loading-overlay').style.display = 'flex';
    }

    hideLoading() {
        document.getElementById('loading-overlay').style.display = 'none';
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
