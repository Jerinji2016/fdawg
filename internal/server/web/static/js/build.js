// Build functionality
class BuildManager {
    constructor() {
        this.platforms = [];
        this.environments = [];
        this.buildConfig = null;
        this.isBuilding = false;
        this.buildProgress = null;
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadInitialData();
    }

    bindEvents() {
        // Refresh buttons
        document.getElementById('refresh-status-btn').addEventListener('click', () => {
            this.loadBuildStatus();
        });

        document.getElementById('refresh-artifacts-btn').addEventListener('click', () => {
            this.loadArtifacts();
        });

        // Setup buttons
        document.getElementById('setup-default-btn').addEventListener('click', () => {
            this.setupDefault();
        });

        document.getElementById('setup-wizard-btn').addEventListener('click', () => {
            this.setupWizard();
        });

        // Config buttons
        document.getElementById('edit-config-btn').addEventListener('click', () => {
            this.editConfig();
        });

        document.getElementById('reset-config-btn').addEventListener('click', () => {
            this.resetConfig();
        });

        // Build buttons
        document.getElementById('start-build-btn').addEventListener('click', () => {
            this.startBuild();
        });

        document.getElementById('stop-build-btn').addEventListener('click', () => {
            this.stopBuild();
        });

        // Toggle sections
        document.getElementById('toggle-build-info').addEventListener('click', () => {
            this.toggleBuildInfoSection();
        });

        // Build Plan Drawer events
        document.getElementById('close-drawer-btn').addEventListener('click', () => {
            this.closeBuildPlanDrawer();
        });

        document.getElementById('cancel-build-btn').addEventListener('click', () => {
            this.closeBuildPlanDrawer();
        });

        document.getElementById('execute-build-btn').addEventListener('click', () => {
            this.executeFromDrawer();
        });

        // Build Configuration Drawer events
        document.getElementById('close-config-drawer-btn').addEventListener('click', () => {
            this.closeBuildConfigDrawer();
        });

        document.getElementById('cancel-config-btn').addEventListener('click', () => {
            this.closeBuildConfigDrawer();
        });

        document.getElementById('save-config-btn').addEventListener('click', () => {
            this.saveConfigFromDrawer();
        });

        // Close drawers when clicking overlay
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('drawer-overlay')) {
                if (e.target.closest('#build-plan-drawer')) {
                    this.closeBuildPlanDrawer();
                } else if (e.target.closest('#build-config-drawer')) {
                    this.closeBuildConfigDrawer();
                }
            }
        });



        // Platform selection change
        document.addEventListener('change', (e) => {
            if (e.target.classList.contains('platform-checkbox')) {
                this.updateBuildButton();
            }
        });
    }

    async loadInitialData() {
        this.showLoading();
        try {
            await Promise.all([
                this.loadPlatforms(),
                this.loadEnvironments(),
                this.loadBuildStatus()
            ]);
        } catch (error) {
            console.error('Error loading initial data:', error);
            showToast('Failed to load build data', 'error');
        } finally {
            this.hideLoading();
        }
    }

    async loadPlatforms() {
        try {
            const response = await fetch('/api/build/platforms');
            if (!response.ok) {
                throw new Error('Failed to load platforms');
            }

            const data = await response.json();
            this.platforms = data.all || [];
            this.renderPlatformSelection();
        } catch (error) {
            console.error('Error loading platforms:', error);
            showToast('Failed to load platforms', 'error');
        }
    }

    async loadEnvironments() {
        try {
            const response = await fetch('/api/environment/list');
            if (!response.ok) {
                console.log('No environments available or API not accessible');
                this.environments = [];
                this.renderEnvironmentSelection();
                return;
            }

            const data = await response.json();
            this.environments = data.environments || data.Environments || [];
            this.renderEnvironmentSelection();
        } catch (error) {
            console.error('Error loading environments:', error);
            this.environments = [];
            this.renderEnvironmentSelection();
        }
    }

    async loadBuildStatus() {
        try {
            const response = await fetch('/api/build/status');
            if (!response.ok) {
                throw new Error('Failed to load build status');
            }

            const data = await response.json();
            this.buildConfig = data.config;
            this.updateStatusDisplay(data);
            this.updateSectionVisibility(data);
        } catch (error) {
            console.error('Error loading build status:', error);
            this.updateStatusDisplay({ 
                config_exists: false, 
                error: error.message 
            });
            this.updateSectionVisibility({ config_exists: false });
        }
    }

    updateStatusDisplay(status) {
        const configStatus = document.getElementById('config-status');
        const lastBuildStatus = document.getElementById('last-build-status');

        if (status.config_exists) {
            configStatus.innerHTML = '<i class="fas fa-check-circle" style="color: green;"></i> Configured';
            lastBuildStatus.innerHTML = status.last_build || '<i class="fas fa-minus"></i> No builds yet';
        } else {
            configStatus.innerHTML = '<i class="fas fa-exclamation-triangle" style="color: orange;"></i> Not configured';
            lastBuildStatus.innerHTML = '<i class="fas fa-minus"></i> Setup required';
        }
    }

    updateSectionVisibility(status) {
        const setupSection = document.getElementById('setup-section');
        const configSection = document.getElementById('config-section');
        const buildSection = document.getElementById('build-section');
        const artifactsSection = document.getElementById('artifacts-section');

        if (status.config_exists) {
            setupSection.style.display = 'none';
            configSection.style.display = 'block';
            buildSection.style.display = 'block';
            artifactsSection.style.display = 'block';
            this.renderConfigDisplay();
            this.loadArtifacts();
        } else {
            setupSection.style.display = 'block';
            configSection.style.display = 'none';
            buildSection.style.display = 'none';
            artifactsSection.style.display = 'none';
        }
    }

    renderPlatformSelection() {
        const container = document.getElementById('platform-selection');

        const html = this.platforms.map(platform => {
            const platformIcon = this.getPlatformIcon(platform.id);
            const isAvailable = platform.available;

            return `
                <label class="platform-checkbox-label ${isAvailable ? '' : 'unavailable'}">
                    <input type="checkbox"
                           class="platform-checkbox"
                           value="${platform.id}"
                           ${isAvailable ? '' : 'disabled'}>
                    <div class="platform-checkbox-content">
                        <i class="${platformIcon}"></i>
                        <span>${platform.name}</span>
                    </div>
                </label>
            `;
        }).join('');

        container.innerHTML = html;
        this.updateBuildButton();
    }

    renderEnvironmentSelection() {
        const select = document.getElementById('environment-select');
        
        // Clear existing options except the first one
        while (select.children.length > 1) {
            select.removeChild(select.lastChild);
        }

        this.environments.forEach(env => {
            const option = document.createElement('option');
            option.value = env.name;
            option.textContent = `${env.name} (${env.description || 'No description'})`;
            select.appendChild(option);
        });
    }

    renderConfigDisplay() {
        const container = document.getElementById('config-display');

        if (!this.buildConfig) {
            container.innerHTML = '<p>Loading configuration...</p>';
            return;
        }



        const html = `
            <div class="config-summary">
                <div class="config-item">
                    <div class="config-label">App Name Source:</div>
                    <div class="config-value">${this.buildConfig.metadata?.app_name_source || this.buildConfig.Metadata?.AppNameSource || 'pubspec'}</div>
                </div>
                <div class="config-item">
                    <div class="config-label">Version Source:</div>
                    <div class="config-value">${this.buildConfig.metadata?.version_source || this.buildConfig.Metadata?.VersionSource || 'pubspec'}</div>
                </div>
                <div class="config-item">
                    <div class="config-label">Output Directory:</div>
                    <div class="config-value">${this.buildConfig.artifacts?.base_output_dir || this.buildConfig.Artifacts?.BaseOutputDir || 'output'}</div>
                </div>
                <div class="config-item">
                    <div class="config-label">Enabled Platforms:</div>
                    <div class="config-value">${this.getEnabledPlatforms()}</div>
                </div>
                <div class="config-item">
                    <div class="config-label">Pre-build Steps:</div>
                    <div class="config-value">${this.getPreBuildSteps()}</div>
                </div>
            </div>
        `;

        container.innerHTML = html;
    }

    getEnabledPlatforms() {
        if (!this.buildConfig) return 'None';

        // Handle both camelCase and PascalCase field names
        const platforms = this.buildConfig.platforms || this.buildConfig.Platforms || {};

        const enabled = [];
        Object.keys(platforms).forEach(platform => {
            const platformConfig = platforms[platform];
            const isEnabled = platformConfig?.enabled || platformConfig?.Enabled;
            if (isEnabled) {
                enabled.push(platform.charAt(0).toUpperCase() + platform.slice(1));
            }
        });

        return enabled.length > 0 ? enabled.join(', ') : 'None configured';
    }

    getPreBuildSteps() {
        if (!this.buildConfig) return 'None';

        // Handle both camelCase and PascalCase field names
        const preBuild = this.buildConfig.pre_build || this.buildConfig.PreBuild || {};

        const steps = [];
        if (preBuild.install_dependencies !== false && preBuild.InstallDependencies !== false) {
            steps.push('Install Dependencies');
        }
        if (preBuild.generate_code || preBuild.GenerateCode) {
            steps.push('Generate Code');
        }
        if (preBuild.custom_commands?.length > 0 || preBuild.CustomCommands?.length > 0) {
            steps.push('Custom Commands');
        }

        return steps.length > 0 ? steps.join(', ') : 'None configured';
    }

    updateBuildButton() {
        const checkboxes = document.querySelectorAll('.platform-checkbox:checked');
        const startBtn = document.getElementById('start-build-btn');
        
        startBtn.disabled = checkboxes.length === 0 || this.isBuilding;
    }

    async setupDefault() {
        this.showConfirmationDialog(
            'Create default build configuration?',
            'This will create a build configuration with default settings for all available platforms.',
            async () => {
                await this.executeSetup(true);
            }
        );
    }

    async setupWizard() {
        this.showConfirmationDialog(
            'Start interactive build setup?',
            'This will guide you through configuring build settings for your project.',
            async () => {
                await this.executeSetup(false);
            }
        );
    }

    async executeSetup(useDefault) {
        this.showLoading();
        
        try {
            const response = await fetch('/api/build/setup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    default: useDefault,
                    force: false
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Setup failed');
            }

            showToast('Build configuration created successfully!', 'success');
            await this.loadBuildStatus();
        } catch (error) {
            console.error('Error during setup:', error);
            showToast(`Setup failed: ${error.message}`, 'error');
        } finally {
            this.hideLoading();
        }
    }

    async startBuild() {
        const selectedPlatforms = Array.from(document.querySelectorAll('.platform-checkbox:checked'))
            .map(cb => cb.value);

        if (selectedPlatforms.length === 0) {
            showToast('Please select at least one platform', 'warning');
            return;
        }

        const environment = document.getElementById('environment-select').value;
        const options = {
            skip_pre_build: document.getElementById('skip-pre-build').checked,
            continue_on_error: document.getElementById('continue-on-error').checked,
            dry_run: document.getElementById('dry-run').checked,
            parallel: document.getElementById('parallel').checked
        };

        const message = `Start build for ${selectedPlatforms.join(', ')}?`;
        const details = `
            Platforms: ${selectedPlatforms.join(', ')}<br>
            Environment: ${environment || 'None'}<br>
            Options: ${Object.keys(options).filter(k => options[k]).join(', ') || 'None'}
        `;

        if (options.dry_run) {
            // For dry-run, directly execute to get the plan and show in drawer
            this.currentBuildParams = { platforms: selectedPlatforms, environment, options };
            await this.executeBuild(selectedPlatforms, environment, options);
        } else {
            this.showConfirmationDialog(message, details, async () => {
                await this.executeBuild(selectedPlatforms, environment, options);
            });
        }
    }

    async executeBuild(platforms, environment, options) {
        this.isBuilding = true;
        this.updateBuildButton();

        const startBtn = document.getElementById('start-build-btn');
        const stopBtn = document.getElementById('stop-build-btn');

        startBtn.style.display = 'none';
        stopBtn.style.display = 'inline-block';

        this.showProgressSection();

        try {
            const response = await fetch('/api/build/run', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    platforms: platforms,
                    environment: environment,
                    ...options
                })
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Build failed');
            }

            const result = await response.json();

            if (result.dry_run) {
                this.showBuildPlanDrawer();
                showToast('Build plan generated!', 'info');
            } else {
                this.displayBuildResult(result);
                showToast('Build completed!', 'success');
            }
        } catch (error) {
            console.error('Build error:', error);
            showToast(`Build failed: ${error.message}`, 'error');
            this.displayBuildError(error.message);
        } finally {
            this.isBuilding = false;
            this.updateBuildButton();
            startBtn.style.display = 'inline-block';
            stopBtn.style.display = 'none';
            this.loadArtifacts(); // Refresh artifacts after build
        }
    }

    async stopBuild() {
        this.showConfirmationDialog(
            'Stop the current build?',
            'This will terminate the build process. Any completed artifacts will be preserved.',
            async () => {
                try {
                    const response = await fetch('/api/build/stop', {
                        method: 'POST'
                    });

                    if (response.ok) {
                        showToast('Build stopped', 'info');
                    }
                } catch (error) {
                    console.error('Error stopping build:', error);
                    showToast('Failed to stop build', 'error');
                }
            }
        );
    }

    showProgressSection() {
        const progressSection = document.getElementById('progress-section');
        const progressContent = document.getElementById('build-progress-content');

        progressSection.style.display = 'block';
        progressContent.innerHTML = `
            <div class="build-progress">
                <div class="progress-header">
                    <i class="fas fa-spinner fa-spin"></i> Build in progress...
                </div>
                <div class="progress-details">
                    <p>Building selected platforms. This may take several minutes.</p>
                    <div class="progress-log" id="build-log">
                        <div class="log-entry">Starting build process...</div>
                    </div>
                </div>
            </div>
        `;
    }

    showBuildPlanDrawer() {
        const drawer = document.getElementById('build-plan-drawer');
        const content = document.getElementById('build-plan-content');

        if (!drawer || !content) {
            return;
        }

        // Store current build parameters for execution
        this.currentBuildParams.options.dry_run = false; // Remove dry-run for actual execution

        const platforms = this.currentBuildParams.platforms;
        const environment = this.currentBuildParams.environment;

        const html = `
            <div class="plan-overview">
                <div class="plan-summary">
                    <h4><i class="fas fa-info-circle"></i> Build Overview</h4>
                    <div class="plan-item">
                        <span class="plan-label">Platforms:</span>
                        <span class="plan-value">${platforms.join(', ')}</span>
                    </div>
                    <div class="plan-item">
                        <span class="plan-label">Environment:</span>
                        <span class="plan-value">${environment || 'None'}</span>
                    </div>
                    <div class="plan-item">
                        <span class="plan-label">Mode:</span>
                        <span class="plan-value">Preview (Dry Run)</span>
                    </div>
                </div>

                <div class="plan-details">
                    <h4><i class="fas fa-cogs"></i> Pre-build Steps</h4>
                    <div class="plan-steps">
                        <div class="plan-step">
                            <i class="fas fa-download"></i>
                            <span>Install dependencies (flutter pub get)</span>
                        </div>
                        ${this.currentBuildParams.options.skip_pre_build ?
                            '<div class="plan-step disabled"><i class="fas fa-times"></i><span>Pre-build steps will be skipped</span></div>' :
                            '<div class="plan-step"><i class="fas fa-code"></i><span>Generate code (if configured)</span></div>'
                        }
                    </div>
                </div>

                <div class="plan-details">
                    <h4><i class="fas fa-hammer"></i> Platform Builds</h4>
                    <div class="plan-platforms">
                        ${platforms.map(platform => `
                            <div class="plan-platform">
                                <i class="${this.getPlatformIcon(platform)}"></i>
                                <div class="platform-details">
                                    <div class="platform-name">${platform.charAt(0).toUpperCase() + platform.slice(1)}</div>
                                    <div class="platform-type">${this.getBuildTypeForPlatform(platform)}</div>
                                </div>
                            </div>
                        `).join('')}
                    </div>
                </div>

                <div class="plan-note">
                    <i class="fas fa-info-circle"></i>
                    <div>
                        <strong>This is a preview.</strong> No files will be created or modified.
                        Click "Execute Build" to run the actual build process.
                    </div>
                </div>
            </div>
        `;

        content.innerHTML = html;
        drawer.classList.add('open');
        document.body.classList.add('drawer-open');
    }

    closeBuildPlanDrawer() {
        const drawer = document.getElementById('build-plan-drawer');
        drawer.classList.remove('open');
        document.body.classList.remove('drawer-open');

        // Reset build state
        this.isBuilding = false;
        this.updateBuildButton();

        const startBtn = document.getElementById('start-build-btn');
        const stopBtn = document.getElementById('stop-build-btn');
        startBtn.style.display = 'inline-block';
        stopBtn.style.display = 'none';
    }

    async executeFromDrawer() {
        if (!this.currentBuildParams) {
            showToast('No build parameters available', 'error');
            return;
        }

        this.closeBuildPlanDrawer();

        // Execute the actual build
        const { platforms, environment, options } = this.currentBuildParams;
        await this.executeBuild(platforms, environment, options);
    }

    getBuildTypeForPlatform(platform) {
        const buildTypes = {
            android: 'APK (Release)',
            ios: 'Archive',
            macos: 'App Bundle',
            linux: 'Executable',
            windows: 'Executable',
            web: 'Web Build'
        };
        return buildTypes[platform] || 'Build';
    }

    displayBuildResult(result) {
        const progressContent = document.getElementById('build-progress-content');

        let html = `
            <div class="build-result">
                <div class="result-header ${result.success ? 'success' : 'error'}">
                    <i class="fas ${result.success ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i>
                    Build ${result.success ? 'Completed' : 'Failed'}
                </div>
                <div class="result-details">
                    <div class="result-item">
                        <span class="result-label">Duration:</span>
                        <span class="result-value">${result.duration || 'Unknown'}</span>
                    </div>
                    <div class="result-item">
                        <span class="result-label">Platforms:</span>
                        <span class="result-value">${Object.keys(result.platform_results || {}).join(', ')}</span>
                    </div>
                </div>
        `;

        if (result.platform_results) {
            html += '<div class="platform-results">';
            Object.entries(result.platform_results).forEach(([platform, platformResult]) => {
                html += `
                    <div class="platform-result ${platformResult.success ? 'success' : 'error'}">
                        <div class="platform-result-header">
                            <i class="${this.getPlatformIcon(platform)}"></i>
                            ${platform.charAt(0).toUpperCase() + platform.slice(1)}
                            <span class="platform-result-status">
                                ${platformResult.success ? 'Success' : 'Failed'}
                            </span>
                        </div>
                        ${platformResult.artifacts ? `
                            <div class="platform-artifacts">
                                Artifacts: ${platformResult.artifacts.length}
                            </div>
                        ` : ''}
                        ${platformResult.error ? `
                            <div class="platform-error">${platformResult.error}</div>
                        ` : ''}
                    </div>
                `;
            });
            html += '</div>';
        }

        html += '</div>';
        progressContent.innerHTML = html;
    }

    displayBuildError(error) {
        const progressContent = document.getElementById('build-progress-content');
        progressContent.innerHTML = `
            <div class="build-result">
                <div class="result-header error">
                    <i class="fas fa-exclamation-circle"></i>
                    Build Failed
                </div>
                <div class="result-details">
                    <div class="error-message">${error}</div>
                </div>
            </div>
        `;
    }

    async loadArtifacts() {
        try {
            const response = await fetch('/api/build/artifacts');
            if (!response.ok) {
                throw new Error('Failed to load artifacts');
            }

            const data = await response.json();
            this.renderArtifacts(data.artifacts || []);
        } catch (error) {
            console.error('Error loading artifacts:', error);
            document.getElementById('artifacts-container').innerHTML =
                '<p class="error-text">Failed to load artifacts</p>';
        }
    }

    renderArtifacts(artifacts) {
        const container = document.getElementById('artifacts-container');

        if (artifacts.length === 0) {
            container.innerHTML = '<p class="info-text">No build artifacts found</p>';
            return;
        }

        const html = artifacts.map(artifact => `
            <div class="artifact-tile">
                <div class="artifact-tile-content">
                    <div class="artifact-tile-header">
                        <div class="artifact-title">
                            <i class="${this.getPlatformIcon(artifact.platform)}"></i>
                            <span class="artifact-name">${artifact.name}</span>
                        </div>
                        <button class="artifact-download-btn" onclick="buildManager.downloadArtifact('${artifact.path}')" title="Download">
                            <i class="fas fa-download"></i>
                        </button>
                    </div>
                    <div class="artifact-tile-body">
                        <div class="artifact-chips">
                            <span class="artifact-chip artifact-type-chip">${artifact.type}</span>
                            <span class="artifact-chip artifact-size-chip">${artifact.size}</span>
                            <span class="artifact-chip artifact-platform-chip">${artifact.platform}</span>
                        </div>
                        <div class="artifact-date">${artifact.date}</div>
                    </div>
                </div>
            </div>
        `).join('');

        container.innerHTML = `<div class="artifacts-grid">${html}</div>`;
    }

    async downloadArtifact(path) {
        try {
            const response = await fetch(`/api/build/artifacts/download?path=${encodeURIComponent(path)}`);
            if (!response.ok) {
                throw new Error('Download failed');
            }

            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = path.split('/').pop();
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        } catch (error) {
            console.error('Download error:', error);
            showToast('Download failed', 'error');
        }
    }

    async editConfig() {
        try {
            const response = await fetch('/api/build/config');
            if (!response.ok) {
                throw new Error('Failed to load configuration');
            }

            const config = await response.json();
            this.showBuildConfigDrawer(config, 'Edit Build Configuration');
        } catch (error) {
            console.error('Error loading config for editing:', error);
            showToast('Failed to load configuration for editing', 'error');
        }
    }

    showBuildConfigDrawer(config, title = 'Build Configuration') {
        const drawer = document.getElementById('build-config-drawer');
        const content = document.getElementById('build-config-content');
        const titleElement = document.getElementById('config-drawer-title');

        if (!drawer || !content) {
            return;
        }

        // Store current config for editing
        this.currentEditConfig = config || {};
        titleElement.textContent = title;

        const html = this.generateConfigDrawerContent(this.currentEditConfig);
        content.innerHTML = html;

        drawer.classList.add('open');
        document.body.classList.add('drawer-open');
    }

    closeBuildConfigDrawer() {
        const drawer = document.getElementById('build-config-drawer');
        drawer.classList.remove('open');
        document.body.classList.remove('drawer-open');
        this.currentEditConfig = null;
    }

    async saveConfigFromDrawer() {
        try {
            const config = this.buildConfigFromDrawerForm();

            const response = await fetch('/api/build/config/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(config)
            });

            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(errorText || 'Failed to save configuration');
            }

            showToast('Configuration saved successfully!', 'success');
            this.closeBuildConfigDrawer();
            await this.loadBuildStatus(); // Refresh the display
        } catch (error) {
            console.error('Error saving config:', error);
            showToast(`Failed to save configuration: ${error.message}`, 'error');
        }
    }

    generateConfigDrawerContent() {
        return `
            <div class="config-drawer-content">
                <!-- Metadata Section -->
                <div class="config-section">
                    <h4><i class="fas fa-info-circle"></i> Metadata Configuration</h4>
                    <div class="config-form-group">
                        <label>App Name Source:</label>
                        <select id="edit-app-name-source" class="config-input">
                            <option value="namer">Namer Configuration</option>
                            <option value="pubspec">pubspec.yaml</option>
                            <option value="custom">Custom</option>
                        </select>
                    </div>
                    <div class="config-form-group" id="custom-app-name-group" style="display: none;">
                        <label>Custom App Name:</label>
                        <input type="text" id="edit-custom-app-name" class="config-input" placeholder="Enter custom app name">
                    </div>
                    <div class="config-form-group">
                        <label>Version Source:</label>
                        <select id="edit-version-source" class="config-input">
                            <option value="pubspec">pubspec.yaml</option>
                            <option value="custom">Custom</option>
                        </select>
                    </div>
                    <div class="config-form-group" id="custom-version-group" style="display: none;">
                        <label>Custom Version:</label>
                        <input type="text" id="edit-custom-version" class="config-input" placeholder="e.g., 1.0.0+1">
                    </div>
                </div>

                <!-- Artifacts Section -->
                <div class="config-section">
                    <h4><i class="fas fa-folder"></i> Artifacts Configuration</h4>
                    <div class="config-form-group">
                        <label>Base Output Directory:</label>
                        <input type="text" id="edit-output-dir" class="config-input" placeholder="build/fdawg-outputs">
                    </div>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-organize-by-date">
                            Organize by date folders
                        </label>
                    </div>
                    <div class="config-form-group" id="date-format-group">
                        <label>Date Format:</label>
                        <select id="edit-date-format" class="config-input">
                            <option value="January-2">January-2 (e.g., June-7)</option>
                            <option value="2006-01-02">2006-01-02 (e.g., 2024-06-07)</option>
                            <option value="02-01-2006">02-01-2006 (e.g., 07-06-2024)</option>
                            <option value="custom">Custom Format</option>
                        </select>
                    </div>
                    <div class="config-form-group" id="custom-date-format-group" style="display: none;">
                        <label>Custom Date Format:</label>
                        <input type="text" id="edit-custom-date-format" class="config-input" placeholder="Go time format (e.g., 2006-01-02)">
                    </div>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-organize-by-platform">
                            Organize by platform folders
                        </label>
                    </div>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-organize-by-build-type">
                            Organize by build type folders
                        </label>
                    </div>
                    <div class="config-form-group">
                        <label>Naming Pattern:</label>
                        <input type="text" id="edit-naming-pattern" class="config-input" placeholder="{app_name}_{version}_{arch}">
                        <small class="config-help">Available variables: {app_name}, {version}, {arch}, {platform}, {build_type}</small>
                    </div>
                    <div class="config-form-group">
                        <label>Fallback App Name:</label>
                        <input type="text" id="edit-fallback-app-name" class="config-input" placeholder="flutter_app">
                    </div>
                </div>

                <!-- Pre-build Steps Section -->
                <div class="config-section">
                    <h4><i class="fas fa-cogs"></i> Pre-build Steps</h4>
                    <div class="config-subsection">
                        <h5><i class="fas fa-globe"></i> Global Steps (All Platforms)</h5>
                        <div id="global-prebuild-steps">
                            <!-- Global pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="global">
                            <i class="fas fa-plus"></i> Add Global Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fab fa-android"></i> Android-specific Steps</h5>
                        <div id="android-prebuild-steps">
                            <!-- Android pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="android">
                            <i class="fas fa-plus"></i> Add Android Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fab fa-apple"></i> iOS-specific Steps</h5>
                        <div id="ios-prebuild-steps">
                            <!-- iOS pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="ios">
                            <i class="fas fa-plus"></i> Add iOS Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fas fa-globe"></i> Web-specific Steps</h5>
                        <div id="web-prebuild-steps">
                            <!-- Web pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="web">
                            <i class="fas fa-plus"></i> Add Web Step
                        </button>
                    </div>
                </div>

                <!-- Platform Configuration Section -->
                <div class="config-section">
                    <h4><i class="fas fa-mobile-alt"></i> Platform Configuration</h4>
                    <div id="platform-config-container">
                        <!-- Platform configurations will be populated here -->
                    </div>
                </div>

                <!-- Execution Configuration Section -->
                <div class="config-section">
                    <h4><i class="fas fa-play"></i> Execution Configuration</h4>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-parallel-builds">
                            Enable parallel builds (experimental)
                        </label>
                    </div>
                    <div class="config-form-group" id="max-parallel-group">
                        <label>Max Parallel Builds:</label>
                        <input type="number" id="edit-max-parallel" class="config-input" min="1" max="8" value="2">
                    </div>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-continue-on-error">
                            Continue on error (don't stop if one platform fails)
                        </label>
                    </div>
                    <div class="config-form-group">
                        <label>
                            <input type="checkbox" id="edit-save-logs">
                            Save build logs
                        </label>
                    </div>
                    <div class="config-form-group">
                        <label>Log Level:</label>
                        <select id="edit-log-level" class="config-input">
                            <option value="debug">Debug</option>
                            <option value="info">Info</option>
                            <option value="warning">Warning</option>
                            <option value="error">Error</option>
                        </select>
                    </div>
                </div>
            </div>
        `;
    }

    buildConfigFromDrawerForm() {
        // Get current config as base
        const config = this.currentEditConfig || {};

        // Update metadata
        config.metadata = {
            app_name_source: document.getElementById('edit-app-name-source').value,
            custom_app_name: document.getElementById('edit-custom-app-name').value,
            version_source: document.getElementById('edit-version-source').value,
            custom_version: document.getElementById('edit-custom-version').value,
        };

        // Update artifacts
        config.artifacts = {
            base_output_dir: document.getElementById('edit-output-dir').value || 'build/fdawg-outputs',
            organization: {
                by_date: document.getElementById('edit-organize-by-date').checked,
                date_format: document.getElementById('edit-date-format').value,
                by_platform: document.getElementById('edit-organize-by-platform').checked,
                by_build_type: document.getElementById('edit-organize-by-build-type').checked,
            },
            naming: {
                pattern: document.getElementById('edit-naming-pattern').value || '{app_name}_{version}_{arch}',
                fallback_app_name: document.getElementById('edit-fallback-app-name').value || 'flutter_app',
            },
        };

        // Update execution
        config.execution = {
            parallel_builds: document.getElementById('edit-parallel-builds').checked,
            max_parallel: parseInt(document.getElementById('edit-max-parallel').value) || 2,
            continue_on_error: document.getElementById('edit-continue-on-error').checked,
            save_logs: document.getElementById('edit-save-logs').checked,
            log_level: document.getElementById('edit-log-level').value,
        };

        // TODO: Add pre-build steps and platform configurations
        // This will be implemented in the next iteration

        return config;
    }

    async setupConfig() {
        // Show configuration drawer with default config for setup
        const defaultConfig = {
            metadata: {
                app_name_source: 'namer',
                version_source: 'pubspec'
            },
            artifacts: {
                base_output_dir: 'build/fdawg-outputs',
                organization: {
                    by_date: true,
                    date_format: 'January-2',
                    by_platform: true,
                    by_build_type: true
                },
                naming: {
                    pattern: '{app_name}_{version}_{arch}',
                    fallback_app_name: 'flutter_app'
                }
            },
            execution: {
                parallel_builds: false,
                max_parallel: 2,
                continue_on_error: false,
                save_logs: true,
                log_level: 'info'
            }
        };

        this.showBuildConfigDrawer(defaultConfig, 'Setup Build Configuration');
    }

    async resetConfig() {
        this.showConfirmationDialog(
            'Reset build configuration?',
            'This will delete the current build configuration. You will need to run setup again.',
            async () => {
                try {
                    const response = await fetch('/api/build/reset', {
                        method: 'POST'
                    });

                    if (!response.ok) {
                        throw new Error('Reset failed');
                    }

                    showToast('Configuration reset successfully', 'success');
                    await this.loadBuildStatus();
                } catch (error) {
                    console.error('Reset error:', error);
                    showToast(`Reset failed: ${error.message}`, 'error');
                }
            }
        );
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

    toggleBuildInfoSection() {
        const container = document.getElementById('build-info-container');
        const toggleBtn = document.getElementById('toggle-build-info');
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

    // Utility methods from namer.js pattern
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

        modal.querySelector('.modal-close').addEventListener('click', () => this.hideConfirmationDialog(modal));
        modal.querySelector('.confirm-no-btn').addEventListener('click', () => this.hideConfirmationDialog(modal));
        modal.querySelector('.confirm-yes-btn').addEventListener('click', async () => {
            this.hideConfirmationDialog(modal);
            await action();
        });

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
}

// Global instance
let buildManager;

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    buildManager = new BuildManager();
});
