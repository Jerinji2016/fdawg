// Build functionality
class BuildManager {
    constructor() {
        this.platforms = [];
        this.environments = [];
        this.buildConfig = null;
        this.isBuilding = false;
        this.buildProgress = null;
        this.buildEventSource = null;
        this.buildResults = {
            succeeded: [],
            failed: [],
            inProgress: []
        };
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
        document.getElementById('collapse-config-btn').addEventListener('click', () => {
            this.toggleConfigPreview();
        });

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
            <div class="build-config-preview" id="build-config-preview-container">
                <div class="preview-content" id="build-config-preview-content">
                    ${this.generateConfigPreview()}
                </div>
            </div>
        `;

        container.innerHTML = html;
    }

    generateConfigPreview() {
        if (!this.buildConfig) return '<p>No configuration available</p>';

        const metadata = this.buildConfig.Metadata || this.buildConfig.metadata || {};
        const artifacts = this.buildConfig.Artifacts || this.buildConfig.artifacts || {};
        const execution = this.buildConfig.Execution || this.buildConfig.execution || {};
        const platforms = this.buildConfig.Platforms || this.buildConfig.platforms || {};
        const preBuild = this.buildConfig.PreBuild || this.buildConfig.pre_build || {};

        return `
            <div class="config-preview-grid">
                <!-- Metadata Section -->
                <div class="config-preview-section">
                    <h5><i class="fas fa-info-circle"></i> Metadata</h5>
                    <div class="config-preview-items">
                        <div class="config-preview-item">
                            <span class="preview-label">App Name Source:</span>
                            <span class="preview-value">${metadata.AppNameSource || metadata.app_name_source || 'namer'}</span>
                        </div>
                        ${metadata.AppNameSource === 'custom' || metadata.app_name_source === 'custom' ? `
                        <div class="config-preview-item">
                            <span class="preview-label">Custom App Name:</span>
                            <span class="preview-value">${metadata.CustomAppName || metadata.custom_app_name || 'Not set'}</span>
                        </div>
                        ` : ''}
                        <div class="config-preview-item">
                            <span class="preview-label">Version Source:</span>
                            <span class="preview-value">${metadata.VersionSource || metadata.version_source || 'pubspec'}</span>
                        </div>
                        ${metadata.VersionSource === 'custom' || metadata.version_source === 'custom' ? `
                        <div class="config-preview-item">
                            <span class="preview-label">Custom Version:</span>
                            <span class="preview-value">${metadata.CustomVersion || metadata.custom_version || 'Not set'}</span>
                        </div>
                        ` : ''}
                    </div>
                </div>

                <!-- Artifacts Section -->
                <div class="config-preview-section">
                    <h5><i class="fas fa-archive"></i> Artifacts</h5>
                    <div class="config-preview-items">
                        <div class="config-preview-item">
                            <span class="preview-label">Output Directory:</span>
                            <span class="preview-value">${artifacts.BaseOutputDir || artifacts.base_output_dir || 'build/fdawg-outputs'}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Organization:</span>
                            <span class="preview-value">${this.getOrganizationSummary(artifacts.Organization || artifacts.organization)}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Naming Pattern:</span>
                            <span class="preview-value">${artifacts.Naming?.Pattern || artifacts.naming?.pattern || '{app_name}_{version}_{arch}'}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Fallback App Name:</span>
                            <span class="preview-value">${artifacts.Naming?.FallbackAppName || artifacts.naming?.fallback_app_name || 'flutter_app'}</span>
                        </div>
                    </div>
                </div>

                <!-- Execution Section -->
                <div class="config-preview-section">
                    <h5><i class="fas fa-cogs"></i> Execution</h5>
                    <div class="config-preview-items">
                        <div class="config-preview-item">
                            <span class="preview-label">Parallel Builds:</span>
                            <span class="preview-value">${execution.ParallelBuilds || execution.parallel_builds ? 'Enabled' : 'Disabled'}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Max Parallel:</span>
                            <span class="preview-value">${execution.MaxParallel || execution.max_parallel || 2}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Continue on Error:</span>
                            <span class="preview-value">${execution.ContinueOnError || execution.continue_on_error ? 'Yes' : 'No'}</span>
                        </div>
                        <div class="config-preview-item">
                            <span class="preview-label">Log Level:</span>
                            <span class="preview-value">${execution.LogLevel || execution.log_level || 'info'}</span>
                        </div>
                    </div>
                </div>

                <!-- Platforms Section -->
                <div class="config-preview-section">
                    <h5><i class="fas fa-mobile-alt"></i> Platforms</h5>
                    <div class="config-preview-items">
                        ${this.getPlatformsSummary(platforms)}
                    </div>
                </div>

                <!-- Pre-build Steps Section -->
                <div class="config-preview-section">
                    <h5><i class="fas fa-list-ol"></i> Pre-build Steps</h5>
                    <div class="config-preview-items">
                        ${this.getPreBuildStepsSummary(preBuild)}
                    </div>
                </div>
            </div>
        `;
    }

    getOrganizationSummary(organization) {
        if (!organization) return 'Default';

        const org = organization;
        const features = [];

        if (org.ByDate || org.by_date) features.push('By Date');
        if (org.ByPlatform || org.by_platform) features.push('By Platform');
        if (org.ByBuildType || org.by_build_type) features.push('By Build Type');

        return features.length > 0 ? features.join(', ') : 'None';
    }

    getPlatformsSummary(platforms) {
        if (!platforms || Object.keys(platforms).length === 0) {
            return '<div class="config-preview-item"><span class="preview-value">No platforms configured</span></div>';
        }

        return Object.keys(platforms).map(platformName => {
            const platform = platforms[platformName];
            const enabled = platform.Enabled !== undefined ? platform.Enabled : platform.enabled;
            const buildTypesCount = platform.BuildTypes?.length || platform.build_types?.length || 0;

            return `
                <div class="config-preview-item">
                    <span class="preview-label">${platformName}:</span>
                    <span class="preview-value ${enabled ? 'enabled' : 'disabled'}">
                        ${enabled ? 'Enabled' : 'Disabled'} (${buildTypesCount} build types)
                    </span>
                </div>
            `;
        }).join('');
    }

    getPreBuildStepsSummary(preBuild) {
        if (!preBuild) return '<div class="config-preview-item"><span class="preview-value">No pre-build steps configured</span></div>';

        const sections = ['Global', 'Android', 'IOS', 'Web'];
        const summaries = [];

        sections.forEach(section => {
            const steps = preBuild[section] || preBuild[section.toLowerCase()] || [];
            if (steps.length > 0) {
                summaries.push(`${section}: ${steps.length} step${steps.length !== 1 ? 's' : ''}`);
            }
        });

        if (summaries.length === 0) {
            return '<div class="config-preview-item"><span class="preview-value">No pre-build steps configured</span></div>';
        }

        return summaries.map(summary => `
            <div class="config-preview-item">
                <span class="preview-value">${summary}</span>
            </div>
        `).join('');
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

    toggleConfigPreview() {
        const configDisplay = document.getElementById('config-display');
        const icon = document.getElementById('config-collapse-icon');

        if (configDisplay.style.display === 'none') {
            configDisplay.style.display = 'block';
            icon.classList.remove('fa-chevron-down');
            icon.classList.add('fa-chevron-up');
        } else {
            configDisplay.style.display = 'none';
            icon.classList.remove('fa-chevron-up');
            icon.classList.add('fa-chevron-down');
        }
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

        // Only show progress section for actual builds, not dry runs
        if (!options.dry_run) {
            this.showProgressSection();

            // Initialize platforms as in progress
            platforms.forEach(platform => {
                this.addPlatformToInProgress(platform);
            });

            this.startBuildStreaming();
        }

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
                // Process build results and update platform status
                this.processBuildResult(result);

                // Update the progress section with final results instead of replacing it
                this.updateProgressSectionWithResults(result);

                // Handle both Go struct naming and JSON naming
                const success = result.Success !== undefined ? result.Success : result.success;

                if (success) {
                    showToast('Build completed!', 'success');
                } else {
                    showToast('Build completed with errors', 'warning');
                }
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
            this.stopBuildStreaming();
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

        // Reset build results
        this.buildResults = {
            succeeded: [],
            failed: [],
            inProgress: []
        };

        progressSection.style.display = 'block';
        progressContent.innerHTML = `
            <div class="build-progress">
                <div class="progress-header">
                    <i class="fas fa-spinner fa-spin"></i> Build in progress...
                </div>
                <div class="progress-details">
                    <p>Building selected platforms. This may take several minutes.</p>
                    <div class="build-status-summary" id="build-status-summary">
                        <div class="status-item">
                            <span class="status-label">In Progress:</span>
                            <span class="status-value" id="in-progress-platforms">-</span>
                        </div>
                        <div class="status-item">
                            <span class="status-label">Succeeded:</span>
                            <span class="status-value succeeded" id="succeeded-platforms">-</span>
                        </div>
                        <div class="status-item">
                            <span class="status-label">Failed:</span>
                            <span class="status-value failed" id="failed-platforms">-</span>
                        </div>
                    </div>
                    <div class="progress-log" id="build-log">
                        <div class="log-entry log-info">
                            <span class="log-timestamp">${new Date().toLocaleTimeString()}</span>
                            <span class="log-icon"><i class="fas fa-info-circle"></i></span>
                            <span class="log-message">Initializing build process...</span>
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    startBuildStreaming() {
        // Close any existing connection
        this.stopBuildStreaming();

        // Create EventSource for streaming build logs
        this.buildEventSource = new EventSource('/api/build/stream');

        this.buildEventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleBuildStreamData(data);
            } catch (error) {
                console.error('Error parsing stream data:', error);
                this.addBuildLogEntry(event.data, 'info');
            }
        };

        this.buildEventSource.onerror = (error) => {
            console.error('Build stream error:', error);
            this.addBuildLogEntry('Connection to build stream lost. Retrying...', 'error');
        };

        this.buildEventSource.onopen = () => {
            console.log('Build stream connected');
            this.addBuildLogEntry('Connected to build stream...', 'info');
        };
    }

    stopBuildStreaming() {
        if (this.buildEventSource) {
            this.buildEventSource.close();
            this.buildEventSource = null;
        }
    }

    handleBuildStreamData(data) {
        switch (data.type) {
            case 'log':
                this.addBuildLogEntry(data.message, data.level || 'info');
                break;
            case 'progress':
                this.updateBuildProgress(data.step, data.total, data.current);
                break;
            case 'status':
                this.updateBuildStatus(data.status, data.message);
                break;
            case 'complete':
                this.handleBuildComplete(data);
                break;
            case 'error':
                this.handleBuildError(data.error);
                break;
            default:
                console.log('Unknown stream data type:', data.type);
        }
    }

    addBuildLogEntry(message, level = 'info') {
        const logContainer = document.getElementById('build-log');
        if (!logContainer) return;

        const timestamp = new Date().toLocaleTimeString();
        const logEntry = document.createElement('div');
        logEntry.className = `log-entry log-${level}`;

        const icon = this.getLogIcon(level);
        logEntry.innerHTML = `
            <span class="log-timestamp">${timestamp}</span>
            <span class="log-icon">${icon}</span>
            <span class="log-message">${this.escapeHtml(message)}</span>
        `;

        logContainer.appendChild(logEntry);
        logContainer.scrollTop = logContainer.scrollHeight;
    }

    getLogIcon(level) {
        const icons = {
            info: '<i class="fas fa-info-circle"></i>',
            success: '<i class="fas fa-check-circle"></i>',
            warning: '<i class="fas fa-exclamation-triangle"></i>',
            error: '<i class="fas fa-times-circle"></i>',
            debug: '<i class="fas fa-bug"></i>'
        };
        return icons[level] || icons.info;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    updateBuildProgress(step, total, current) {
        const progressHeader = document.querySelector('.progress-header');
        if (progressHeader) {
            progressHeader.innerHTML = `
                <i class="fas fa-spinner fa-spin"></i>
                Building... (${current}/${total}) - ${step}
            `;
        }
    }

    updateBuildStatus(status, message) {
        const progressHeader = document.querySelector('.progress-header');
        if (progressHeader) {
            const icon = status === 'success' ? 'fa-check' :
                        status === 'error' ? 'fa-times' : 'fa-spinner fa-spin';
            progressHeader.innerHTML = `
                <i class="fas ${icon}"></i> ${message}
            `;
        }
    }

    handleBuildComplete() {
        this.addBuildLogEntry('Build completed successfully!', 'success');
        this.updateBuildStatus('success', 'Build completed!');
        this.stopBuildStreaming();
        showToast('Build completed successfully!', 'success');
    }

    handleBuildError(error) {
        this.addBuildLogEntry(`Build failed: ${error}`, 'error');
        this.updateBuildStatus('error', 'Build failed');
        this.stopBuildStreaming();
        showToast(`Build failed: ${error}`, 'error');
    }

    updateBuildStatusSummary() {
        const inProgressEl = document.getElementById('in-progress-platforms');
        const succeededEl = document.getElementById('succeeded-platforms');
        const failedEl = document.getElementById('failed-platforms');

        if (inProgressEl) {
            inProgressEl.textContent = this.buildResults.inProgress.length > 0 ?
                this.buildResults.inProgress.join(', ') : '-';
        }

        if (succeededEl) {
            succeededEl.textContent = this.buildResults.succeeded.length > 0 ?
                this.buildResults.succeeded.join(', ') : '-';
        }

        if (failedEl) {
            failedEl.textContent = this.buildResults.failed.length > 0 ?
                this.buildResults.failed.join(', ') : '-';
        }
    }

    addPlatformToInProgress(platform) {
        if (!this.buildResults.inProgress.includes(platform)) {
            this.buildResults.inProgress.push(platform);
            this.updateBuildStatusSummary();
        }
    }

    movePlatformToSucceeded(platform) {
        // Remove from in progress
        const inProgressIndex = this.buildResults.inProgress.indexOf(platform);
        if (inProgressIndex > -1) {
            this.buildResults.inProgress.splice(inProgressIndex, 1);
        }

        // Add to succeeded if not already there
        if (!this.buildResults.succeeded.includes(platform)) {
            this.buildResults.succeeded.push(platform);
        }

        this.updateBuildStatusSummary();
        this.addBuildLogEntry(`✅ ${platform} build completed successfully`, 'success');
    }

    movePlatformToFailed(platform, error) {
        // Remove from in progress
        const inProgressIndex = this.buildResults.inProgress.indexOf(platform);
        if (inProgressIndex > -1) {
            this.buildResults.inProgress.splice(inProgressIndex, 1);
        }

        // Add to failed if not already there
        if (!this.buildResults.failed.includes(platform)) {
            this.buildResults.failed.push(platform);
        }

        this.updateBuildStatusSummary();
        this.addBuildLogEntry(`❌ ${platform} build failed: ${error}`, 'error');
    }

    processBuildResult(result) {
        console.log('Processing build result:', result);

        // Clear in-progress platforms since build is complete
        this.buildResults.inProgress = [];

        // Process platform results if available (handle both Go naming conventions)
        const platformResults = result.PlatformResults || result.platform_results;
        console.log('Platform results:', platformResults);

        if (platformResults) {
            Object.entries(platformResults).forEach(([platform, platformResult]) => {
                // Handle both Go struct naming (Success) and JSON naming (success)
                const success = platformResult.Success !== undefined ? platformResult.Success : platformResult.success;
                const error = platformResult.Error || platformResult.error;

                console.log(`Platform ${platform}: success=${success}, error=${error}`);

                if (success) {
                    this.movePlatformToSucceeded(platform);
                } else {
                    this.movePlatformToFailed(platform, error ? error.toString() : 'Build failed');
                }
            });
        } else {
            // Fallback: if no platform results, assume all platforms succeeded if overall success
            const platforms = this.currentBuildParams?.platforms || [];
            platforms.forEach(platform => {
                // Handle both Go struct naming (Success) and JSON naming (success)
                const success = result.Success !== undefined ? result.Success : result.success;
                console.log(`Fallback for platform ${platform}: success=${success}`);

                if (success) {
                    this.movePlatformToSucceeded(platform);
                } else {
                    this.movePlatformToFailed(platform, 'Build failed');
                }
            });
        }

        // Update final status
        this.updateBuildStatusSummary();

        // Update progress header
        const totalPlatforms = this.buildResults.succeeded.length + this.buildResults.failed.length;
        const successCount = this.buildResults.succeeded.length;

        if (this.buildResults.failed.length === 0) {
            this.updateBuildStatus('success', `All ${totalPlatforms} platform(s) built successfully!`);
        } else if (this.buildResults.succeeded.length === 0) {
            this.updateBuildStatus('error', `All ${totalPlatforms} platform(s) failed to build`);
        } else {
            this.updateBuildStatus('warning', `${successCount}/${totalPlatforms} platform(s) built successfully`);
        }
    }

    updateProgressSectionWithResults(result) {
        const progressContent = document.getElementById('build-progress-content');
        if (!progressContent) return;

        // Handle both Go struct naming and JSON naming
        const success = result.Success !== undefined ? result.Success : result.success;
        const duration = result.Duration || result.duration || 'Unknown';
        const platformResults = result.PlatformResults || result.platform_results || {};

        // Keep the existing build progress structure but add results
        const existingProgress = progressContent.querySelector('.build-progress');
        if (existingProgress) {
            // Add build results section to existing progress
            const resultsSection = document.createElement('div');
            resultsSection.className = 'build-results-section';
            resultsSection.innerHTML = `
                <div class="results-header">
                    <h4><i class="fas ${success ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i> Build Results</h4>
                </div>
                <div class="results-details">
                    <div class="result-item">
                        <span class="result-label">Duration:</span>
                        <span class="result-value">${this.formatDuration(duration)}</span>
                    </div>
                    <div class="result-item">
                        <span class="result-label">Total Artifacts:</span>
                        <span class="result-value">${(result.Artifacts || result.artifacts || []).length}</span>
                    </div>
                </div>
                ${Object.keys(platformResults).length > 0 ? this.generatePlatformResultsHTML(platformResults) : ''}
            `;

            existingProgress.appendChild(resultsSection);
        }
    }

    generatePlatformResultsHTML(platformResults) {
        let html = '<div class="platform-results-summary"><h5>Platform Details</h5>';

        Object.entries(platformResults).forEach(([platform, platformResult]) => {
            const platformSuccess = platformResult.Success !== undefined ? platformResult.Success : platformResult.success;
            const artifacts = platformResult.Artifacts || platformResult.artifacts || [];
            const error = platformResult.Error || platformResult.error;

            html += `
                <div class="platform-result-item ${platformSuccess ? 'success' : 'error'}">
                    <div class="platform-result-header">
                        <i class="${this.getPlatformIcon(platform)}"></i>
                        <span class="platform-name">${platform.charAt(0).toUpperCase() + platform.slice(1)}</span>
                        <span class="platform-status ${platformSuccess ? 'success' : 'error'}">
                            ${platformSuccess ? '✅ Success' : '❌ Failed'}
                        </span>
                    </div>
                    ${artifacts.length > 0 ? `
                        <div class="platform-artifacts-count">
                            ${artifacts.length} artifact${artifacts.length !== 1 ? 's' : ''} created
                        </div>
                    ` : ''}
                    ${error ? `
                        <div class="platform-error-msg">${error.toString ? error.toString() : error}</div>
                    ` : ''}
                </div>
            `;
        });

        html += '</div>';
        return html;
    }

    formatDuration(duration) {
        if (typeof duration === 'string') return duration;
        if (typeof duration === 'number') {
            // Convert nanoseconds to seconds
            const seconds = duration / 1000000000;
            return `${seconds.toFixed(2)}s`;
        }
        return 'Unknown';
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

        // Handle both Go struct naming and JSON naming
        const success = result.Success !== undefined ? result.Success : result.success;
        const duration = result.Duration || result.duration || 'Unknown';
        const platformResults = result.PlatformResults || result.platform_results || {};

        let html = `
            <div class="build-result">
                <div class="result-header ${success ? 'success' : 'error'}">
                    <i class="fas ${success ? 'fa-check-circle' : 'fa-exclamation-circle'}"></i>
                    Build ${success ? 'Completed' : 'Failed'}
                </div>
                <div class="result-details">
                    <div class="result-item">
                        <span class="result-label">Duration:</span>
                        <span class="result-value">${duration}</span>
                    </div>
                    <div class="result-item">
                        <span class="result-label">Platforms:</span>
                        <span class="result-value">${Object.keys(platformResults).join(', ')}</span>
                    </div>
                </div>
        `;

        if (Object.keys(platformResults).length > 0) {
            html += '<div class="platform-results">';
            Object.entries(platformResults).forEach(([platform, platformResult]) => {
                // Handle both Go struct naming and JSON naming
                const platformSuccess = platformResult.Success !== undefined ? platformResult.Success : platformResult.success;
                const platformError = platformResult.Error || platformResult.error;
                const artifacts = platformResult.Artifacts || platformResult.artifacts || [];

                html += `
                    <div class="platform-result ${platformSuccess ? 'success' : 'error'}">
                        <div class="platform-result-header">
                            <i class="${this.getPlatformIcon(platform)}"></i>
                            ${platform.charAt(0).toUpperCase() + platform.slice(1)}
                            <span class="platform-result-status">
                                ${platformSuccess ? 'Success' : 'Failed'}
                            </span>
                        </div>
                        ${artifacts.length > 0 ? `
                            <div class="platform-artifacts">
                                Artifacts: ${artifacts.length}
                            </div>
                        ` : ''}
                        ${platformError ? `
                            <div class="platform-error">${platformError.toString ? platformError.toString() : platformError}</div>
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

        // Populate form with current config
        this.populateConfigDrawerForm(this.currentEditConfig);

        // Populate platform configuration
        this.populatePlatformConfig(this.currentEditConfig);

        // Bind dynamic events
        this.bindConfigDrawerEvents();

        // Bind add step buttons after HTML is inserted
        this.bindAddStepButtons();

        drawer.classList.add('open');
        document.body.classList.add('drawer-open');
    }

    closeBuildConfigDrawer() {
        const drawer = document.getElementById('build-config-drawer');
        drawer.classList.remove('open');
        document.body.classList.remove('drawer-open');
        this.currentEditConfig = null;
    }

    validateConfigForm() {
        const errors = [];

        // Validate app name source
        const appNameSource = document.getElementById('edit-app-name-source')?.value;
        if (!appNameSource) {
            errors.push('App name source is required');
        }

        // Validate custom app name if custom source is selected
        if (appNameSource === 'custom') {
            const customAppName = document.getElementById('edit-custom-app-name')?.value;
            if (!customAppName || customAppName.trim() === '') {
                errors.push('Custom app name is required when using custom source');
            }
        }

        // Validate version source
        const versionSource = document.getElementById('edit-version-source')?.value;
        if (!versionSource) {
            errors.push('Version source is required');
        }

        // Validate custom version if custom source is selected
        if (versionSource === 'custom') {
            const customVersion = document.getElementById('edit-custom-version')?.value;
            if (!customVersion || customVersion.trim() === '') {
                errors.push('Custom version is required when using custom source');
            }
        }

        // Validate output directory
        const outputDir = document.getElementById('edit-output-dir')?.value;
        if (!outputDir || outputDir.trim() === '') {
            errors.push('Base output directory is required');
        }

        // Validate naming pattern
        const namingPattern = document.getElementById('edit-naming-pattern')?.value;
        if (!namingPattern || namingPattern.trim() === '') {
            errors.push('Naming pattern is required');
        }

        // Validate fallback app name
        const fallbackAppName = document.getElementById('edit-fallback-app-name')?.value;
        if (!fallbackAppName || fallbackAppName.trim() === '') {
            errors.push('Fallback app name is required');
        }

        if (errors.length > 0) {
            showToast(`Validation errors: ${errors.join(', ')}`, 'error');
            return false;
        }

        return true;
    }

    populatePlatformConfig(config) {
        const container = document.getElementById('platform-config-container');
        if (!container) return;

        const platforms = config?.Platforms || config?.platforms || {};
        const platformNames = ['Android', 'IOS', 'Web', 'MacOS', 'Linux', 'Windows'];

        let html = '<div class="platform-config-grid">';

        platformNames.forEach(platformName => {
            const platformConfig = platforms[platformName] || platforms[platformName.toLowerCase()] || {};
            const enabled = platformConfig.Enabled !== undefined ? platformConfig.Enabled : true;
            const buildTypesCount = platformConfig.BuildTypes ? platformConfig.BuildTypes.length : 0;

            html += `
                <div class="platform-config-card">
                    <div class="platform-config-header">
                        <div class="platform-info">
                            <i class="${this.getPlatformIcon(platformName.toLowerCase())}"></i>
                            <span class="platform-name">${platformName}</span>
                        </div>
                        <label class="toggle-switch">
                            <input type="checkbox" class="platform-enabled-toggle"
                                   data-platform="${platformName}" ${enabled ? 'checked' : ''}>
                            <span class="toggle-slider"></span>
                        </label>
                    </div>
                    <div class="platform-config-details">
                        <span class="build-types-count">${buildTypesCount} build type${buildTypesCount !== 1 ? 's' : ''}</span>
                        <button type="button" class="secondary-btn small-btn configure-platform-btn"
                                data-platform="${platformName}">
                            <i class="fas fa-cog"></i> Configure
                        </button>
                    </div>
                </div>
            `;
        });

        html += '</div>';
        container.innerHTML = html;

        // Bind platform configuration events
        this.bindPlatformConfigEvents();
    }

    bindPlatformConfigEvents() {
        // Platform enabled/disabled toggles
        document.querySelectorAll('.platform-enabled-toggle').forEach(toggle => {
            toggle.addEventListener('change', (e) => {
                const platform = e.target.dataset.platform;
                const enabled = e.target.checked;
                console.log(`Platform ${platform} ${enabled ? 'enabled' : 'disabled'}`);
                // Update the current config
                if (!this.currentEditConfig.Platforms) {
                    this.currentEditConfig.Platforms = {};
                }
                if (!this.currentEditConfig.Platforms[platform]) {
                    this.currentEditConfig.Platforms[platform] = {};
                }
                this.currentEditConfig.Platforms[platform].Enabled = enabled;
            });
        });

        // Configure platform buttons
        document.querySelectorAll('.configure-platform-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const platform = e.target.dataset.platform || e.target.closest('.configure-platform-btn').dataset.platform;
                this.showPlatformConfigModal(platform);
            });
        });
    }

    showPlatformConfigModal(platform) {
        showToast(`Platform configuration for ${platform} coming soon!`, 'info');
        // TODO: Implement detailed platform configuration modal
    }

    getPlatformIcon(platform) {
        const icons = {
            'android': 'fab fa-android',
            'ios': 'fab fa-apple',
            'web': 'fas fa-globe',
            'macos': 'fab fa-apple',
            'linux': 'fab fa-linux',
            'windows': 'fab fa-windows'
        };
        return icons[platform.toLowerCase()] || 'fas fa-desktop';
    }

    async saveConfigFromDrawer() {
        try {
            // Validate form before building config
            if (!this.validateConfigForm()) {
                return;
            }

            const config = this.buildConfigFromDrawerForm();

            console.log('Saving config:', config);

            const response = await fetch('/api/build/config/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(config)
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Server error response:', errorText);
                throw new Error(errorText || 'Failed to save configuration');
            }

            const result = await response.json();
            console.log('Save result:', result);

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
                    <p class="config-description">Configure how app name and version are determined for builds. Choose between automatic detection from project files or manual specification.</p>
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
                    <p class="config-description">Control where build artifacts are stored and how they are organized. Configure output directories, folder structure, and file naming patterns.</p>
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
                    <p class="config-description">Define commands that run before building. Global steps apply to all platforms, while platform-specific steps only run for those targets.</p>
                    <div class="config-subsection">
                        <h5><i class="fas fa-globe"></i> Global Steps (All Platforms)</h5>
                        <div id="global-steps-container">
                            <!-- Global pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="global">
                            <i class="fas fa-plus"></i> Add Global Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fab fa-android"></i> Android-specific Steps</h5>
                        <div id="android-steps-container">
                            <!-- Android pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="android">
                            <i class="fas fa-plus"></i> Add Android Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fab fa-apple"></i> iOS-specific Steps</h5>
                        <div id="ios-steps-container">
                            <!-- iOS pre-build steps will be populated here -->
                        </div>
                        <button type="button" class="secondary-btn add-step-btn" data-platform="ios">
                            <i class="fas fa-plus"></i> Add iOS Step
                        </button>
                    </div>
                    <div class="config-subsection">
                        <h5><i class="fas fa-globe"></i> Web-specific Steps</h5>
                        <div id="web-steps-container">
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
                    <p class="config-description">Enable or disable specific platforms for building. Only enabled platforms will be available for selection during builds.</p>
                    <div id="platform-config-container">
                        <!-- Platform configurations will be populated here -->
                    </div>
                </div>

                <!-- Execution Configuration Section -->
                <div class="config-section">
                    <h4><i class="fas fa-play"></i> Execution Configuration</h4>
                    <p class="config-description">Control how builds are executed, including parallel processing, error handling, and logging options for better build management.</p>
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

    populateConfigDrawerForm(config) {
        if (!config) return;

        console.log('Populating form with config:', config);

        // Metadata
        const metadata = config.metadata || config.Metadata || {};
        const appNameSourceEl = document.getElementById('edit-app-name-source');
        const customAppNameEl = document.getElementById('edit-custom-app-name');
        const versionSourceEl = document.getElementById('edit-version-source');
        const customVersionEl = document.getElementById('edit-custom-version');

        if (appNameSourceEl) appNameSourceEl.value = metadata.app_name_source || metadata.AppNameSource || 'namer';
        if (customAppNameEl) customAppNameEl.value = metadata.custom_app_name || metadata.CustomAppName || '';
        if (versionSourceEl) versionSourceEl.value = metadata.version_source || metadata.VersionSource || 'pubspec';
        if (customVersionEl) customVersionEl.value = metadata.custom_version || metadata.CustomVersion || '';

        // Artifacts
        const artifacts = config.artifacts || config.Artifacts || {};
        const organization = artifacts.organization || artifacts.Organization || {};
        const naming = artifacts.naming || artifacts.Naming || {};

        const outputDirEl = document.getElementById('edit-output-dir');
        const organizeDateEl = document.getElementById('edit-organize-by-date');
        const dateFormatEl = document.getElementById('edit-date-format');
        const organizePlatformEl = document.getElementById('edit-organize-by-platform');
        const organizeBuildTypeEl = document.getElementById('edit-organize-by-build-type');
        const namingPatternEl = document.getElementById('edit-naming-pattern');
        const fallbackAppNameEl = document.getElementById('edit-fallback-app-name');

        if (outputDirEl) outputDirEl.value = artifacts.base_output_dir || artifacts.BaseOutputDir || 'build/fdawg-outputs';
        if (organizeDateEl) organizeDateEl.checked = organization.by_date || organization.ByDate || false;
        if (dateFormatEl) dateFormatEl.value = organization.date_format || organization.DateFormat || 'January-2';
        if (organizePlatformEl) organizePlatformEl.checked = organization.by_platform || organization.ByPlatform || false;
        if (organizeBuildTypeEl) organizeBuildTypeEl.checked = organization.by_build_type || organization.ByBuildType || false;
        if (namingPatternEl) namingPatternEl.value = naming.pattern || naming.Pattern || '{app_name}_{version}_{arch}';
        if (fallbackAppNameEl) fallbackAppNameEl.value = naming.fallback_app_name || naming.FallbackAppName || 'flutter_app';

        // Execution
        const execution = config.execution || config.Execution || {};
        const parallelBuildsEl = document.getElementById('edit-parallel-builds');
        const maxParallelEl = document.getElementById('edit-max-parallel');
        const continueOnErrorEl = document.getElementById('edit-continue-on-error');
        const saveLogsEl = document.getElementById('edit-save-logs');
        const logLevelEl = document.getElementById('edit-log-level');

        if (parallelBuildsEl) parallelBuildsEl.checked = execution.parallel_builds || execution.ParallelBuilds || false;
        if (maxParallelEl) maxParallelEl.value = execution.max_parallel || execution.MaxParallel || 2;
        if (continueOnErrorEl) continueOnErrorEl.checked = execution.continue_on_error || execution.ContinueOnError || false;
        if (saveLogsEl) saveLogsEl.checked = execution.save_logs !== undefined ? (execution.save_logs || execution.SaveLogs) : true;
        if (logLevelEl) logLevelEl.value = execution.log_level || execution.LogLevel || 'info';

        // Update visibility of custom fields
        this.updateCustomFieldVisibility();

        // Populate prebuild steps
        this.populatePreBuildSteps(config);
    }

    populatePreBuildSteps(config) {
        if (!config) return;

        const preBuild = config.PreBuild || config.pre_build || {};
        const platforms = ['global', 'android', 'ios', 'web'];

        platforms.forEach(platform => {
            const steps = preBuild[platform] || preBuild[platform.charAt(0).toUpperCase() + platform.slice(1)] || [];
            const container = document.getElementById(`${platform}-steps-container`);

            if (!container) return;

            // Clear existing steps
            container.innerHTML = '';

            // Add existing steps
            steps.forEach((step) => {
                this.addPreBuildStep(platform, step);
            });

            // If no steps exist, ensure there's at least one empty step
            if (steps.length === 0) {
                this.addPreBuildStep(platform);
            }
        });
    }

    bindConfigDrawerEvents() {
        // App name source change
        document.getElementById('edit-app-name-source').addEventListener('change', () => {
            this.updateCustomFieldVisibility();
        });

        // Version source change
        document.getElementById('edit-version-source').addEventListener('change', () => {
            this.updateCustomFieldVisibility();
        });

        // Date format change
        document.getElementById('edit-date-format').addEventListener('change', () => {
            this.updateCustomFieldVisibility();
        });

        // Parallel builds change
        document.getElementById('edit-parallel-builds').addEventListener('change', () => {
            this.updateCustomFieldVisibility();
        });

        // Add step buttons - use event delegation on the drawer
        const drawer = document.getElementById('build-config-drawer');
        if (drawer) {
            // Remove any existing listeners to prevent duplicates
            drawer.removeEventListener('click', this.handleAddStepClick);

            // Add new listener
            this.handleAddStepClick = (e) => {
                if (e.target.classList.contains('add-step-btn') || e.target.closest('.add-step-btn')) {
                    const btn = e.target.classList.contains('add-step-btn') ? e.target : e.target.closest('.add-step-btn');
                    const platform = btn.dataset.platform;
                    if (platform) {
                        this.addPreBuildStep(platform);
                    }
                }
            };

            drawer.addEventListener('click', this.handleAddStepClick);
        }
    }

    bindAddStepButtons() {
        const addStepButtons = document.querySelectorAll('.add-step-btn');

        addStepButtons.forEach(btn => {
            const platform = btn.dataset.platform;

            // Add new listener
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.stopPropagation();
                this.addPreBuildStep(platform);
            });
        });
    }

    updateCustomFieldVisibility() {
        // App name custom field
        const appNameSource = document.getElementById('edit-app-name-source').value;
        const customAppNameGroup = document.getElementById('custom-app-name-group');
        if (customAppNameGroup) {
            customAppNameGroup.style.display = appNameSource === 'custom' ? 'block' : 'none';
        }

        // Version custom field
        const versionSource = document.getElementById('edit-version-source').value;
        const customVersionGroup = document.getElementById('custom-version-group');
        if (customVersionGroup) {
            customVersionGroup.style.display = versionSource === 'custom' ? 'block' : 'none';
        }

        // Date format custom field
        const dateFormat = document.getElementById('edit-date-format').value;
        const customDateFormatGroup = document.getElementById('custom-date-format-group');
        if (customDateFormatGroup) {
            customDateFormatGroup.style.display = dateFormat === 'custom' ? 'block' : 'none';
        }

        // Max parallel field
        const parallelBuilds = document.getElementById('edit-parallel-builds').checked;
        const maxParallelGroup = document.getElementById('max-parallel-group');
        if (maxParallelGroup) {
            maxParallelGroup.style.display = parallelBuilds ? 'block' : 'none';
        }
    }

    addPreBuildStep(platform, existingStep = null) {
        const container = document.getElementById(`${platform}-steps-container`);
        if (!container) {
            return;
        }

        const stepIndex = container.children.length;

        // Use existing step data or defaults
        const stepName = existingStep?.Name || existingStep?.name || `Custom Step ${stepIndex + 1}`;
        const stepCommand = existingStep?.Command || existingStep?.command || '';
        const stepRequired = existingStep?.Required !== undefined ? existingStep.Required :
                           (existingStep?.required !== undefined ? existingStep.required : true);
        const stepTimeout = existingStep?.Timeout || existingStep?.timeout || 300;

        const stepHTML = `
            <div class="prebuild-step-item" data-platform="${platform}" data-index="${stepIndex}">
                <div class="step-header">
                    <input type="text" class="config-input step-name" placeholder="Step name" value="${stepName}">
                    <button type="button" class="remove-step-btn" title="Remove step">
                        <i class="fas fa-times"></i>
                    </button>
                </div>
                <div class="step-body">
                    <div class="config-form-group">
                        <label>Command:</label>
                        <input type="text" class="config-input step-command" placeholder="e.g., dart run build_runner build" value="${stepCommand}">
                    </div>
                    <div class="step-options">
                        <div class="step-option-group checkbox-group">
                            <label class="checkbox-label">
                                <input type="checkbox" class="step-required" ${stepRequired ? 'checked' : ''}>
                                <span class="checkbox-text">Required</span>
                            </label>
                            <div class="info-tooltip">
                                <i class="fas fa-info-circle info-icon"></i>
                                <span class="tooltip-text">Fail build if this step fails</span>
                            </div>
                        </div>
                        <div class="step-option-group">
                            <label class="input-label">Timeout (seconds):</label>
                            <input type="number" class="config-input step-timeout" value="${stepTimeout}" min="30" max="3600" placeholder="300">
                        </div>
                    </div>
                </div>
            </div>
        `;

        container.insertAdjacentHTML('beforeend', stepHTML);

        // Bind remove button
        const removeBtn = container.querySelector('.prebuild-step-item:last-child .remove-step-btn');
        removeBtn.addEventListener('click', () => {
            removeBtn.closest('.prebuild-step-item').remove();
        });
    }

    buildConfigFromDrawerForm() {
        // Create a fresh config object with only snake_case fields
        const config = {};

        // Update metadata
        const appNameSourceEl = document.getElementById('edit-app-name-source');
        const versionSourceEl = document.getElementById('edit-version-source');

        if (!appNameSourceEl || !versionSourceEl) {
            console.error('Required form elements not found');
            throw new Error('Required form elements not found');
        }

        const appNameSource = appNameSourceEl.value || 'namer';
        const versionSource = versionSourceEl.value || 'pubspec';

        config.Metadata = {
            AppNameSource: appNameSource,
            CustomAppName: '',
            VersionSource: versionSource,
            CustomVersion: '',
        };

        // Add custom values only if custom source is selected
        if (appNameSource === 'custom') {
            const customAppNameEl = document.getElementById('edit-custom-app-name');
            config.Metadata.CustomAppName = customAppNameEl ? customAppNameEl.value : '';
        }
        if (versionSource === 'custom') {
            const customVersionEl = document.getElementById('edit-custom-version');
            config.Metadata.CustomVersion = customVersionEl ? customVersionEl.value : '';
        }

        // Update artifacts
        const dateFormatEl = document.getElementById('edit-date-format');
        const dateFormat = dateFormatEl ? dateFormatEl.value || 'January-2' : 'January-2';

        const customDateFormatEl = document.getElementById('edit-custom-date-format');
        const finalDateFormat = dateFormat === 'custom' && customDateFormatEl ?
            customDateFormatEl.value || 'January-2' : dateFormat;

        const organizationConfig = {
            ByDate: document.getElementById('edit-organize-by-date')?.checked || false,
            DateFormat: finalDateFormat,
            ByPlatform: document.getElementById('edit-organize-by-platform')?.checked || false,
            ByBuildType: document.getElementById('edit-organize-by-build-type')?.checked || false,
        };

        config.Artifacts = {
            BaseOutputDir: document.getElementById('edit-output-dir')?.value || 'build/fdawg-outputs',
            Organization: organizationConfig,
            Naming: {
                Pattern: document.getElementById('edit-naming-pattern')?.value || '{app_name}_{version}_{arch}',
                FallbackAppName: document.getElementById('edit-fallback-app-name')?.value || 'flutter_app',
            },
            Cleanup: {
                Enabled: true,
                KeepLastBuilds: 10,
                MaxAgeDays: 30,
            },
        };

        // Update execution
        const maxParallelEl = document.getElementById('edit-max-parallel');
        const logLevelEl = document.getElementById('edit-log-level');

        config.Execution = {
            ParallelBuilds: document.getElementById('edit-parallel-builds')?.checked || false,
            MaxParallel: maxParallelEl ? parseInt(maxParallelEl.value) || 2 : 2,
            ContinueOnError: document.getElementById('edit-continue-on-error')?.checked || false,
            SaveLogs: document.getElementById('edit-save-logs')?.checked || true,
            LogLevel: logLevelEl ? logLevelEl.value || 'info' : 'info',
        };

        // Update pre-build steps
        config.PreBuild = {
            Global: this.collectPreBuildSteps('global'),
            Android: this.collectPreBuildSteps('android'),
            IOS: this.collectPreBuildSteps('ios'),
            Web: this.collectPreBuildSteps('web'),
        };

        // Update platforms (preserve existing platform configurations from currentEditConfig)
        const existingConfig = this.currentEditConfig || {};
        config.Platforms = existingConfig.Platforms || existingConfig.platforms || {};

        // Ensure platforms config has proper structure if empty
        if (!config.Platforms || Object.keys(config.Platforms).length === 0) {
            config.Platforms = {
                Android: { Enabled: true, BuildTypes: [], Environment: {} },
                IOS: { Enabled: true, BuildTypes: [] },
                Web: { Enabled: true, BuildTypes: [] },
                MacOS: { Enabled: true, BuildTypes: [] },
                Linux: { Enabled: true, BuildTypes: [] },
                Windows: { Enabled: true, BuildTypes: [] }
            };
        }

        console.log('Built config from form:', config);
        return config;
    }

    collectPreBuildSteps(platform) {
        const container = document.getElementById(`${platform}-steps-container`);
        if (!container) return [];

        const steps = [];
        const stepItems = container.querySelectorAll('.prebuild-step-item');

        stepItems.forEach(item => {
            const name = item.querySelector('.step-name').value.trim();
            const command = item.querySelector('.step-command').value.trim();
            const required = item.querySelector('.step-required').checked;
            const timeout = parseInt(item.querySelector('.step-timeout').value) || 300;

            if (name && command) {
                steps.push({
                    Name: name,
                    Command: command,
                    Required: required,
                    Timeout: timeout,
                    WorkingDir: '',
                    Environment: {},
                    Condition: ''
                });
            }
        });

        return steps;
    }

    async setupConfig() {
        // Show configuration drawer with default config for setup
        const defaultConfig = {
            Metadata: {
                AppNameSource: 'namer',
                VersionSource: 'pubspec',
                CustomAppName: '',
                CustomVersion: ''
            },
            Artifacts: {
                BaseOutputDir: 'build/fdawg-outputs',
                Organization: {
                    ByDate: true,
                    DateFormat: 'January-2',
                    ByPlatform: true,
                    ByBuildType: true
                },
                Naming: {
                    Pattern: '{app_name}_{version}_{arch}',
                    FallbackAppName: 'flutter_app'
                },
                Cleanup: {
                    Enabled: true,
                    KeepLastBuilds: 10,
                    MaxAgeDays: 30
                }
            },
            Execution: {
                ParallelBuilds: false,
                MaxParallel: 2,
                ContinueOnError: false,
                SaveLogs: true,
                LogLevel: 'info'
            },
            PreBuild: {
                Global: [],
                Android: [],
                IOS: [],
                Web: []
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
