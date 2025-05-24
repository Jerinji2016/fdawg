// Localization management JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log("Localization management loaded");

    // Initialize localization data
    let localizationData = {
        languages: [],
        translationKeys: [],
        stats: {
            supportedLanguages: 0,
            translationKeys: 0,
            missingTranslations: 0,
            completionRate: 0
        }
    };

    // Translation configuration
    let translationConfig = {
        enabled: false,
        hasApiKey: false
    };

    // Toggle localization summary section
    const toggleSummaryBtn = document.getElementById('toggle-localization-summary');
    const summaryContent = document.getElementById('localization-summary-content');

    if (toggleSummaryBtn && summaryContent) {
        toggleSummaryBtn.addEventListener('click', function() {
            summaryContent.classList.toggle('collapsed');
            toggleSummaryBtn.classList.toggle('collapsed');
        });
    }

    // Add language button event listener
    const addLanguageBtn = document.getElementById('add-language-btn');
    if (addLanguageBtn) {
        addLanguageBtn.addEventListener('click', function() {
            showAddLanguageModal();
        });
    }

    // Add translation key button event listener
    const addTranslationKeyBtn = document.getElementById('add-translation-key-btn');
    if (addTranslationKeyBtn) {
        addTranslationKeyBtn.addEventListener('click', function() {
            showAddTranslationKeyModal();
        });
    }

    // Search functionality
    const searchInput = document.getElementById('translation-search');
    if (searchInput) {
        searchInput.addEventListener('input', function() {
            filterTranslationKeys(this.value);
        });
    }

    // Load initial data
    loadLocalizationData();
    loadTranslationConfig();
    updateConfigurationUI();

    // Function to load localization data from the server
    function loadLocalizationData() {
        // Show loading indicators
        showLoadingIndicators();

        // Fetch localization data from API
        fetch('/api/localizations/data')
            .then(response => response.json())
            .then(data => {
                localizationData = data;
                updateUI();
                hideLoadingIndicators();
            })
            .catch(error => {
                console.error('Error loading localization data:', error);
                hideLoadingIndicators();
                showErrorToast('Failed to load localization data', 'Error');
            });
    }

    // Function to load translation configuration
    function loadTranslationConfig() {
        fetch('/api/localizations/translate-config')
            .then(response => response.json())
            .then(data => {
                translationConfig = data;
                console.log('Translation config loaded:', translationConfig);
                updateConfigurationUI();
            })
            .catch(error => {
                console.error('Error loading translation config:', error);
                translationConfig = { enabled: false, hasApiKey: false };
                updateConfigurationUI();
            });
    }

    // Function to show loading indicators
    function showLoadingIndicators() {
        const loadingLanguages = document.getElementById('loading-languages');
        const loadingTranslations = document.getElementById('loading-translations');

        if (loadingLanguages) loadingLanguages.style.display = 'block';
        if (loadingTranslations) loadingTranslations.style.display = 'block';
    }

    // Function to hide loading indicators
    function hideLoadingIndicators() {
        const loadingLanguages = document.getElementById('loading-languages');
        const loadingTranslations = document.getElementById('loading-translations');

        if (loadingLanguages) loadingLanguages.style.display = 'none';
        if (loadingTranslations) loadingTranslations.style.display = 'none';
    }

    // Function to update the UI with loaded data
    function updateUI() {
        updateSummaryCards();
        updateLanguageCards();
        updateTranslationTable();
    }

    // Function to update summary cards
    function updateSummaryCards() {
        const stats = localizationData.stats;

        document.getElementById('supported-languages-count').textContent = stats.supportedLanguages;
        document.getElementById('translation-keys-count').textContent = stats.translationKeys;
        document.getElementById('missing-translations-count').textContent = stats.missingTranslations;
        document.getElementById('completion-rate').textContent = stats.completionRate + '%';
    }

    // Function to update language cards
    function updateLanguageCards() {
        const languageCardsContainer = document.getElementById('language-cards');

        // Sort languages alphabetically by name
        const languages = [...localizationData.languages].sort((a, b) => a.name.localeCompare(b.name));

        if (languages.length === 0) {
            languageCardsContainer.innerHTML = `
                <div class="no-data-message">
                    <p>No languages configured. Add your first language to get started.</p>
                </div>
            `;
            return;
        }

        let cardsHTML = '';
        languages.forEach(language => {
            const statusClass = getLanguageStatusClass(language.completionRate);
            const statusText = getLanguageStatusText(language.completionRate, language.missingKeys);

            cardsHTML += `
                <div class="language-card" data-language="${language.code}">
                    <div class="language-card-header">
                        <div class="language-flag">${language.flag}</div>
                        <div class="language-actions">
                            <button class="icon-btn download-btn" title="Download ${language.name}" data-language="${language.code}">
                                <i class="fas fa-download"></i>
                            </button>
                            <button class="icon-btn delete-language-btn" title="Delete ${language.name}" data-language="${language.code}">
                                <i class="fas fa-trash"></i>
                            </button>
                        </div>
                    </div>
                    <div class="language-info">
                        <span class="language-name">${language.name}</span>
                        <span class="language-code">${language.code}</span>
                    </div>
                    <div class="language-status">
                        <span class="status-badge ${statusClass}">${statusText}</span>
                    </div>
                </div>
            `;
        });

        languageCardsContainer.innerHTML = cardsHTML;

        // Add event listeners to the new buttons
        addLanguageCardEventListeners();
    }

    // Function to add event listeners to language card buttons
    function addLanguageCardEventListeners() {
        // Download buttons
        document.querySelectorAll('.download-btn').forEach(btn => {
            btn.addEventListener('click', function(e) {
                e.stopPropagation();
                const languageCode = this.getAttribute('data-language');
                downloadLanguageFile(languageCode);
            });
        });

        // Delete buttons
        document.querySelectorAll('.delete-language-btn').forEach(btn => {
            btn.addEventListener('click', function(e) {
                e.stopPropagation();
                const languageCode = this.getAttribute('data-language');
                showDeleteLanguageConfirmation(languageCode);
            });
        });
    }

    // Function to get language status class
    function getLanguageStatusClass(completionRate) {
        if (completionRate === 100) return 'complete';
        if (completionRate === 0) return 'empty';
        return 'incomplete';
    }

    // Function to get language status text
    function getLanguageStatusText(completionRate, missingKeys) {
        if (completionRate === 100) return 'Complete';
        if (completionRate === 0) return 'Empty';
        return `Missing ${missingKeys} keys`;
    }

    // Function to update translation table
    function updateTranslationTable() {
        const table = document.getElementById('translation-table');
        const tbody = document.getElementById('translation-table-body');
        const noDataMessage = document.getElementById('no-translations-message');

        // Sort languages alphabetically by name
        const languages = [...localizationData.languages].sort((a, b) => a.name.localeCompare(b.name));

        // Sort translation keys alphabetically by key
        const translationKeys = [...localizationData.translationKeys].sort((a, b) => a.key.localeCompare(b.key));

        // Update table headers
        updateTableHeaders(table, languages);

        if (translationKeys.length === 0) {
            tbody.innerHTML = '';
            noDataMessage.style.display = 'block';
            return;
        }

        noDataMessage.style.display = 'none';

        // Build table rows
        let rowsHTML = '';
        translationKeys.forEach(keyData => {
            rowsHTML += `
                <tr data-key="${keyData.key}">
                    <td class="editable-key" data-original-key="${keyData.key}" title="Double-click to edit key">
                        <span class="key-text">${keyData.key}</span>
                        <button class="expand-btn" style="display: none;">...</button>
                        <div class="edit-container" style="display: none;">
                            <textarea class="key-input" rows="3">${keyData.key}</textarea>
                            <button class="table-btn save-key-btn" title="Save key">
                                <i class="fas fa-check"></i>
                            </button>
                        </div>
                    </td>
            `;

            // Add language columns
            languages.forEach(language => {
                const value = keyData.translations[language.code] || '';
                const isEmpty = !value.trim();
                rowsHTML += `
                    <td class="editable-translation ${isEmpty ? 'empty-translation' : ''}"
                        data-key="${keyData.key}"
                        data-language="${language.code}"
                        title="Double-click to edit translation">
                        <span class="translation-text">${isEmpty ? '<em>Missing</em>' : value}</span>
                        <button class="expand-btn" style="display: none;">...</button>
                        <button class="table-btn translate-btn" style="display: none;" title="Translate to ${language.name}">
                            <i class="fas fa-language"></i>
                        </button>
                        <div class="edit-container" style="display: none;">
                            <textarea class="translation-input" rows="3">${value}</textarea>
                            <button class="table-btn save-translation-btn" title="Save translation">
                                <i class="fas fa-check"></i>
                            </button>
                        </div>
                    </td>
                `;
            });

            rowsHTML += `
                    <td>
                        <button class="table-btn translate-row-btn" data-key="${keyData.key}" title="Translate row">
                            <i class="fas fa-language"></i>
                        </button>
                        <button class="table-btn delete-key-btn" data-key="${keyData.key}" title="Delete key">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
        });

        tbody.innerHTML = rowsHTML;

        // Add event listeners to table elements
        addTableEventListeners();

        // Update translate button states
        updateTranslateButtonStates();

        // Check for text overflow and add expand buttons (with delay to ensure DOM is rendered)
        setTimeout(() => {
            checkTextOverflow();
        }, 100);
    }

    // Function to update table headers
    function updateTableHeaders(table, languages) {
        const thead = table.querySelector('thead tr');
        const actionsHeader = thead.querySelector('th:last-child');

        // Remove existing language headers
        const existingHeaders = thead.querySelectorAll('th:not(:first-child):not(:last-child)');
        existingHeaders.forEach(header => header.remove());

        // Add language headers
        languages.forEach(language => {
            const th = document.createElement('th');
            th.textContent = language.name;
            thead.insertBefore(th, actionsHeader);
        });
    }

    // Function to add event listeners to table elements
    function addTableEventListeners() {
        // Delete buttons
        document.querySelectorAll('.delete-key-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                const key = this.getAttribute('data-key');
                showDeleteKeyConfirmation(key);
            });
        });

        // Editable key cells
        document.querySelectorAll('.editable-key').forEach(cell => {
            const keyText = cell.querySelector('.key-text');
            const editContainer = cell.querySelector('.edit-container');
            const keyInput = cell.querySelector('.key-input');
            const saveBtn = cell.querySelector('.save-key-btn');
            const expandBtn = cell.querySelector('.expand-btn');

            // Double-click to edit key (whole cell clickable)
            cell.addEventListener('dblclick', function(e) {
                // Prevent editing if already in edit mode or clicking on buttons
                if (cell.classList.contains('editing') ||
                    e.target.classList.contains('save-key-btn') ||
                    e.target.classList.contains('expand-btn')) {
                    return;
                }
                enterEditMode(cell, keyText, editContainer, keyInput, saveBtn);
            });

            // Save key button
            saveBtn.addEventListener('click', function() {
                const originalKey = cell.getAttribute('data-original-key');
                const newKey = keyInput.value.trim();
                if (newKey && newKey !== originalKey) {
                    updateTranslationKey(originalKey, newKey, cell, keyText, editContainer, keyInput, saveBtn);
                } else {
                    exitEditMode(cell, keyText, editContainer, keyInput, saveBtn, originalKey);
                }
            });

            // Enter key to save
            keyInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    saveBtn.click();
                }
            });

            // Escape key to cancel
            keyInput.addEventListener('keydown', function(e) {
                if (e.key === 'Escape') {
                    const originalKey = cell.getAttribute('data-original-key');
                    exitEditMode(cell, keyText, editContainer, keyInput, saveBtn, originalKey);
                }
            });

            // Expand/collapse button
            if (expandBtn) {
                expandBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    toggleRowExpansion(cell);
                });
            }
        });

        // Editable translation cells
        document.querySelectorAll('.editable-translation').forEach(cell => {
            const translationText = cell.querySelector('.translation-text');
            const editContainer = cell.querySelector('.edit-container');
            const translationInput = cell.querySelector('.translation-input');
            const saveBtn = cell.querySelector('.save-translation-btn');
            const expandBtn = cell.querySelector('.expand-btn');

            // Double-click to edit translation (whole cell clickable)
            cell.addEventListener('dblclick', function(e) {
                // Prevent editing if already in edit mode or clicking on buttons
                if (cell.classList.contains('editing') ||
                    e.target.classList.contains('save-translation-btn') ||
                    e.target.classList.contains('expand-btn')) {
                    return;
                }
                enterEditMode(cell, translationText, editContainer, translationInput, saveBtn);
            });

            // Save translation button
            saveBtn.addEventListener('click', function() {
                const key = cell.getAttribute('data-key');
                const language = cell.getAttribute('data-language');
                const newValue = translationInput.value.trim();
                updateSingleTranslation(key, language, newValue, cell, translationText, editContainer, translationInput, saveBtn);
            });

            // Enter key to save
            translationInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    saveBtn.click();
                }
            });

            // Escape key to cancel
            translationInput.addEventListener('keydown', function(e) {
                if (e.key === 'Escape') {
                    const originalValue = translationInput.getAttribute('data-original-value') || '';
                    exitEditMode(cell, translationText, editContainer, translationInput, saveBtn, originalValue);
                }
            });

            // Expand/collapse button
            if (expandBtn) {
                expandBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    toggleRowExpansion(cell);
                });
            }

            // Translate button
            const translateBtn = cell.querySelector('.translate-btn');
            if (translateBtn) {
                translateBtn.addEventListener('click', function(e) {
                    e.stopPropagation();
                    const key = cell.getAttribute('data-key');
                    const targetLanguage = cell.getAttribute('data-language');
                    handleCellTranslation(key, targetLanguage);
                });
            }

            // Add hover functionality for translate button
            cell.addEventListener('mouseenter', function() {
                if (!cell.classList.contains('editing') && translateBtn && translationConfig.enabled) {
                    translateBtn.style.display = 'block';
                }
            });

            cell.addEventListener('mouseleave', function() {
                if (translateBtn) {
                    translateBtn.style.display = 'none';
                }
            });
        });

        // Row translate button handlers
        document.querySelectorAll('.translate-row-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                if (!this.disabled && !this.classList.contains('disabled')) {
                    const key = this.getAttribute('data-key');
                    handleRowTranslation(key);
                }
            });
        });
    }

    // Function to filter translation keys
    function filterTranslationKeys(searchTerm) {
        const rows = document.querySelectorAll('#translation-table-body tr');
        const term = searchTerm.toLowerCase();

        rows.forEach(row => {
            const key = row.getAttribute('data-key').toLowerCase();
            const shouldShow = key.includes(term);
            row.style.display = shouldShow ? '' : 'none';
        });
    }

    // Helper function to enter edit mode
    function enterEditMode(cell, textElement, editContainer, inputElement, saveBtn) {
        // Store original value for cancellation
        let originalValue;
        if (textElement.innerHTML.includes('<em>Missing</em>')) {
            originalValue = '';
        } else {
            originalValue = textElement.textContent.trim();
        }
        inputElement.setAttribute('data-original-value', originalValue);

        // Hide text and expand button, show edit container
        textElement.style.display = 'none';
        const expandBtn = cell.querySelector('.expand-btn');
        if (expandBtn) {
            expandBtn.style.display = 'none';
        }
        editContainer.style.display = 'flex';

        // Focus and select input
        setTimeout(() => {
            inputElement.focus();
            inputElement.select();
        }, 10);

        // Add editing class for styling
        cell.classList.add('editing');
    }

    // Helper function to exit edit mode
    function exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, displayValue) {
        // Show text, hide edit container
        textElement.style.display = 'block';
        editContainer.style.display = 'none';

        // Show expand button if needed
        const expandBtn = cell.querySelector('.expand-btn');
        if (expandBtn && cell.classList.contains('has-overflow')) {
            expandBtn.style.display = 'block';
        }

        // Update display value
        if (displayValue !== undefined) {
            if (displayValue.trim() === '') {
                textElement.innerHTML = '<em>Missing</em>';
                cell.classList.add('empty-translation');
            } else {
                textElement.textContent = displayValue;
                cell.classList.remove('empty-translation');
            }
            inputElement.value = displayValue;
        }

        // Remove editing class
        cell.classList.remove('editing');

        // Recheck text overflow after update
        checkCellOverflow(cell);
    }

    // Function to check text overflow and show expand buttons
    function checkTextOverflow() {
        document.querySelectorAll('.editable-key, .editable-translation').forEach(cell => {
            checkCellOverflow(cell);
        });
    }

    // Function to check if a single cell has text overflow
    function checkCellOverflow(cell) {
        const textElement = cell.querySelector('.key-text, .translation-text');
        const expandBtn = cell.querySelector('.expand-btn');

        if (!textElement || !expandBtn) return;

        // Skip if text is empty or just "Missing"
        const textContent = textElement.textContent.trim();
        if (!textContent || textContent === 'Missing') {
            cell.classList.remove('has-overflow');
            expandBtn.style.display = 'none';
            return;
        }

        // Create a temporary element to measure text height
        const tempElement = textElement.cloneNode(true);
        tempElement.style.position = 'absolute';
        tempElement.style.visibility = 'hidden';
        tempElement.style.height = 'auto';
        tempElement.style.maxHeight = 'none';
        tempElement.style.webkitLineClamp = 'unset';
        tempElement.style.overflow = 'visible';
        tempElement.style.width = textElement.offsetWidth + 'px';

        document.body.appendChild(tempElement);
        const fullHeight = tempElement.offsetHeight;
        document.body.removeChild(tempElement);

        // Get the height when clamped to 2 lines
        const lineHeight = parseFloat(getComputedStyle(textElement).lineHeight) || 20;
        const maxHeight = lineHeight * 2;

        // Show expand button if text exceeds 2 lines
        if (fullHeight > maxHeight + 5) { // Add small tolerance
            cell.classList.add('has-overflow');
            expandBtn.style.display = 'block';
            expandBtn.textContent = textElement.classList.contains('expanded') ? '−' : '...';
        } else {
            cell.classList.remove('has-overflow');
            expandBtn.style.display = 'none';
        }
    }

    // Function to toggle expansion for entire row
    function toggleRowExpansion(clickedCell) {
        const row = clickedCell.closest('tr');
        if (!row) return;

        // Find all cells in this row that have expandable content
        const expandableCells = row.querySelectorAll('.editable-key, .editable-translation');
        const isCurrentlyExpanded = clickedCell.querySelector('.key-text, .translation-text').classList.contains('expanded');

        // Toggle expansion for all cells in the row
        expandableCells.forEach(cell => {
            const textElement = cell.querySelector('.key-text, .translation-text');
            const expandBtn = cell.querySelector('.expand-btn');

            if (!textElement || !expandBtn) return;

            if (isCurrentlyExpanded) {
                // Collapse all
                textElement.classList.remove('expanded');
                expandBtn.textContent = '...';
            } else {
                // Expand all
                textElement.classList.add('expanded');
                expandBtn.textContent = '−';
            }
        });
    }

    // Function to update a single translation
    function updateSingleTranslation(key, language, newValue, cell, textElement, editContainer, inputElement, saveBtn) {
        const formData = new FormData();
        formData.append('translation_key', key);

        // Get existing translations for this key to preserve them
        const keyData = localizationData.translationKeys.find(k => k.key === key);
        const existingTranslations = keyData ? { ...keyData.translations } : {};

        // Update only the specific language
        existingTranslations[language] = newValue;

        formData.append('translations', JSON.stringify(existingTranslations));

        fetch('/api/localizations/update-translations', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Update the local data for this specific language only
                const keyData = localizationData.translationKeys.find(k => k.key === key);
                if (keyData) {
                    keyData.translations[language] = newValue;
                }

                // Exit edit mode with new value
                exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, newValue);

                showSuccessToast(`Translation for "${language}" updated successfully`, 'Success');
            } else {
                showErrorToast(data.error || 'Failed to update translation', 'Error');
                // Exit edit mode with original value
                const originalValue = inputElement.getAttribute('data-original-value') || '';
                exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, originalValue);
            }
        })
        .catch(error => {
            console.error('Error updating translation:', error);
            showErrorToast('Failed to update translation', 'Error');
            // Exit edit mode with original value
            const originalValue = inputElement.getAttribute('data-original-value') || '';
            exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, originalValue);
        });
    }

    // Function to update a translation key
    function updateTranslationKey(originalKey, newKey, cell, textElement, editContainer, inputElement, saveBtn) {
        // First, get all current translations for this key
        const keyData = localizationData.translationKeys.find(k => k.key === originalKey);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Error');
            exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, originalKey);
            return;
        }

        // Delete the old key first
        const deleteFormData = new FormData();
        deleteFormData.append('translation_key', originalKey);

        fetch('/api/localizations/delete-key', {
            method: 'POST',
            body: deleteFormData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Now add the new key with all the translations
                const addFormData = new FormData();
                addFormData.append('translation_key', newKey);
                addFormData.append('translations', JSON.stringify(keyData.translations));

                return fetch('/api/localizations/update-translations', {
                    method: 'POST',
                    body: addFormData
                });
            } else {
                throw new Error(data.error || 'Failed to delete old key');
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Update the row data-key attribute
                const row = cell.closest('tr');
                row.setAttribute('data-key', newKey);
                cell.setAttribute('data-original-key', newKey);

                // Update all translation cells in this row
                const translationCells = row.querySelectorAll('.editable-translation');
                translationCells.forEach(translationCell => {
                    translationCell.setAttribute('data-key', newKey);
                });

                // Update delete button
                const deleteBtn = row.querySelector('.delete-key-btn');
                if (deleteBtn) {
                    deleteBtn.setAttribute('data-key', newKey);
                }

                // Exit edit mode with new key
                exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, newKey);

                showSuccessToast(`Translation key updated successfully`, 'Success');

                // Reload data to ensure consistency
                loadLocalizationData();
            } else {
                throw new Error(data.error || 'Failed to create new key');
            }
        })
        .catch(error => {
            console.error('Error updating translation key:', error);
            showErrorToast('Failed to update translation key', 'Error');
            // Exit edit mode with original key
            exitEditMode(cell, textElement, editContainer, inputElement, saveBtn, originalKey);
        });
    }

    // Modal and action functions
    function showAddLanguageModal() {
        showInputDialog(
            'Add New Language',
            'Enter language code (e.g., de, it, pt_BR):',
            '',
            'Language Code',
            (languageCode) => {
                if (languageCode && languageCode.trim()) {
                    addLanguage(languageCode.trim());
                }
            }
        );
    }

    function showAddTranslationKeyModal() {
        showInputDialog(
            'Add Translation Key',
            'Enter translation key (e.g., auth.logout, common.cancel):',
            '',
            'Translation Key',
            (key) => {
                if (key && key.trim()) {
                    addTranslationKey(key.trim());
                }
            }
        );
    }



    function showDeleteLanguageConfirmation(languageCode) {
        const language = localizationData.languages.find(lang => lang.code === languageCode);
        const languageName = language ? language.name : languageCode;

        showConfirmationToast(
            `This will permanently delete all translations for ${languageName} and cannot be undone.`,
            `Delete ${languageName}?`,
            {
                confirmText: 'Delete',
                cancelText: 'Cancel',
                confirmButtonClass: 'primary-btn',
                onConfirm: () => {
                    deleteLanguage(languageCode);
                }
            }
        );
    }

    function showDeleteKeyConfirmation(key) {
        showConfirmationToast(
            `This will permanently delete the translation key "${key}" from all languages and cannot be undone.`,
            'Delete Translation Key?',
            {
                confirmText: 'Delete',
                cancelText: 'Cancel',
                confirmButtonClass: 'primary-btn',
                onConfirm: () => {
                    deleteTranslationKey(key);
                }
            }
        );
    }

    function downloadLanguageFile(languageCode) {
        const language = localizationData.languages.find(lang => lang.code === languageCode);
        const languageName = language ? language.name : languageCode;

        // Create download link
        const downloadUrl = `/api/localizations/download/${languageCode}`;
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = `${languageCode}.json`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);

        showSuccessToast(`Downloaded ${languageName} translation file`, 'Download Complete');
    }

    // API functions for language management
    function addLanguage(languageCode) {
        const formData = new FormData();
        formData.append('language_code', languageCode);

        fetch('/api/localizations/add-language', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showSuccessToast(`Language ${languageCode} added successfully`, 'Success');
                loadLocalizationData(); // Reload data
            } else {
                showErrorToast(data.error || 'Failed to add language', 'Error');
            }
        })
        .catch(error => {
            console.error('Error adding language:', error);
            showErrorToast('Failed to add language', 'Error');
        });
    }

    function deleteLanguage(languageCode) {
        const formData = new FormData();
        formData.append('language_code', languageCode);

        fetch('/api/localizations/delete-language', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showSuccessToast(`Language deleted successfully`, 'Success');
                loadLocalizationData(); // Reload data
            } else {
                showErrorToast(data.error || 'Failed to delete language', 'Error');
            }
        })
        .catch(error => {
            console.error('Error deleting language:', error);
            showErrorToast('Failed to delete language', 'Error');
        });
    }

    function addTranslationKey(key) {
        const formData = new FormData();
        formData.append('translation_key', key);

        fetch('/api/localizations/add-key', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showSuccessToast(`Translation key "${key}" added successfully`, 'Success');
                loadLocalizationData(); // Reload data
            } else {
                showErrorToast(data.error || 'Failed to add translation key', 'Error');
            }
        })
        .catch(error => {
            console.error('Error adding translation key:', error);
            showErrorToast('Failed to add translation key', 'Error');
        });
    }

    function deleteTranslationKey(key) {
        const formData = new FormData();
        formData.append('translation_key', key);

        fetch('/api/localizations/delete-key', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showSuccessToast(`Translation key "${key}" deleted successfully`, 'Success');
                loadLocalizationData(); // Reload data
            } else {
                showErrorToast(data.error || 'Failed to delete translation key', 'Error');
            }
        })
        .catch(error => {
            console.error('Error deleting translation key:', error);
            showErrorToast('Failed to delete translation key', 'Error');
        });
    }

    // Translation functions
    function updateTranslateButtonStates() {
        if (!translationConfig.enabled) {
            // Hide all translate buttons if translation is disabled
            document.querySelectorAll('.translate-btn, .translate-row-btn').forEach(btn => {
                btn.style.display = 'none';
            });
            return;
        }

        // Update row translate button states
        document.querySelectorAll('.translate-row-btn').forEach(btn => {
            const key = btn.getAttribute('data-key');
            const keyData = localizationData.translationKeys.find(k => k.key === key);

            if (keyData && hasTranslatableContent(keyData.translations)) {
                btn.classList.remove('disabled');
                btn.disabled = false;
                const availableLanguages = getAvailableSourceLanguages(keyData.translations);
                btn.title = `Translate from ${availableLanguages.length} available language(s)`;
            } else {
                btn.classList.add('disabled');
                btn.disabled = true;
                btn.title = 'No content available to translate from';
            }
        });
    }

    function handleCellTranslation(key, targetLanguage) {
        if (!translationConfig.enabled) {
            showErrorToast('Translation service is not enabled. Please configure Google Translate API key in the Web UI.', 'Translation Error');
            return;
        }

        const keyData = localizationData.translationKeys.find(k => k.key === key);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Translation Error');
            return;
        }

        const availableSourceLanguages = getAvailableSourceLanguages(keyData.translations, targetLanguage);

        if (availableSourceLanguages.length === 0) {
            showErrorToast('No content available to translate from', 'Translation Error');
            return;
        }

        if (availableSourceLanguages.length === 1) {
            // Auto-translate from the only available source
            translateCell(key, targetLanguage, availableSourceLanguages[0]);
        } else {
            // Show selection modal with only languages that have content
            showSourceLanguageModal(availableSourceLanguages, (sourceLanguage) => {
                translateCell(key, targetLanguage, sourceLanguage);
            });
        }
    }

    function handleRowTranslation(key) {
        if (!translationConfig.enabled) {
            showErrorToast('Translation service is not enabled. Please configure Google Translate API key in the Web UI.', 'Translation Error');
            return;
        }

        const keyData = localizationData.translationKeys.find(k => k.key === key);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Translation Error');
            return;
        }

        const availableSourceLanguages = getAvailableSourceLanguages(keyData.translations);

        if (availableSourceLanguages.length === 0) {
            showErrorToast('No content available to translate from', 'Translation Error');
            return;
        }

        if (availableSourceLanguages.length === 1) {
            // Auto-translate from the only available source
            const sourceLanguage = availableSourceLanguages[0];
            const targetLanguages = localizationData.languages.map(lang => lang.code);
            translateRow(key, sourceLanguage, targetLanguages);
        } else {
            // Show selection modal with only languages that have content
            showSourceLanguageModal(availableSourceLanguages, (sourceLanguage) => {
                const targetLanguages = localizationData.languages.map(lang => lang.code);
                translateRow(key, sourceLanguage, targetLanguages);
            });
        }
    }

    function translateCell(key, targetLanguage, sourceLanguage) {
        const keyData = localizationData.translationKeys.find(k => k.key === key);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Translation Error');
            return;
        }

        // Find the cell element
        const cell = document.querySelector(`[data-key="${key}"][data-language="${targetLanguage}"]`);
        if (!cell) {
            showErrorToast('Translation cell not found', 'Translation Error');
            return;
        }

        // Show loading indicator
        const translationText = cell.querySelector('.translation-text');
        const originalContent = translationText.innerHTML;
        translationText.innerHTML = 'Translating... <span class="translating-spinner"></span>';

        const formData = new FormData();
        formData.append('translation_key', key);
        formData.append('target_language', targetLanguage);
        formData.append('source_language', sourceLanguage);
        formData.append('existing_translations', JSON.stringify(keyData.translations));

        fetch('/api/localizations/translate-cell', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success && data.translated_text) {
                // Update the translation in the UI
                translationText.textContent = data.translated_text;
                cell.classList.remove('empty-translation');

                // Update local data
                keyData.translations[targetLanguage] = data.translated_text;

                // Save the translation to the server
                updateSingleTranslationFromAPI(key, targetLanguage, data.translated_text);

                showSuccessToast(`Translated to ${getLanguageName(targetLanguage)}`, 'Translation Complete');
            } else {
                translationText.innerHTML = originalContent;
                showErrorToast(data.error || 'Translation failed', 'Translation Error');
            }
        })
        .catch(error => {
            console.error('Error translating cell:', error);
            translationText.innerHTML = originalContent;
            showErrorToast('Translation failed', 'Translation Error');
        });
    }

    function translateRow(key, sourceLanguage, targetLanguages) {
        const keyData = localizationData.translationKeys.find(k => k.key === key);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Translation Error');
            return;
        }

        // Show loading indicators for all target cells
        const loadingCells = new Map();
        targetLanguages.forEach(targetLang => {
            if (targetLang !== sourceLanguage) {
                const cell = document.querySelector(`[data-key="${key}"][data-language="${targetLang}"]`);
                if (cell) {
                    const translationText = cell.querySelector('.translation-text');
                    const isEmpty = !keyData.translations[targetLang] || keyData.translations[targetLang].trim() === '';

                    if (isEmpty) {
                        const originalContent = translationText.innerHTML;
                        loadingCells.set(targetLang, { cell, translationText, originalContent });
                        translationText.innerHTML = 'Translating... <span class="translating-spinner"></span>';
                    }
                }
            }
        });

        const formData = new FormData();
        formData.append('translation_key', key);
        formData.append('source_language', sourceLanguage);
        formData.append('target_languages', JSON.stringify(targetLanguages));
        formData.append('existing_translations', JSON.stringify(keyData.translations));

        fetch('/api/localizations/translate-row', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.translations) {
                let successCount = 0;
                let errorCount = 0;

                // Update each translation
                Object.entries(data.translations).forEach(([targetLang, result]) => {
                    const cellData = loadingCells.get(targetLang);
                    if (cellData) {
                        const { cell, translationText } = cellData;

                        if (result.success && result.text) {
                            translationText.textContent = result.text;
                            cell.classList.remove('empty-translation');
                            keyData.translations[targetLang] = result.text;
                            successCount++;
                        } else {
                            translationText.innerHTML = cellData.originalContent;
                            errorCount++;
                        }
                    }
                });

                // Save all successful translations
                if (successCount > 0) {
                    updateSingleTranslationFromAPI(key, null, null, keyData.translations);
                }

                if (successCount > 0 && errorCount === 0) {
                    showSuccessToast(`Translated ${successCount} language(s)`, 'Translation Complete');
                } else if (successCount > 0 && errorCount > 0) {
                    showSuccessToast(`Translated ${successCount} language(s), ${errorCount} failed`, 'Partial Success');
                } else {
                    showErrorToast('All translations failed', 'Translation Error');
                }
            } else {
                // Restore original content for all loading cells
                loadingCells.forEach(({ translationText, originalContent }) => {
                    translationText.innerHTML = originalContent;
                });
                showErrorToast(data.error || 'Translation failed', 'Translation Error');
            }
        })
        .catch(error => {
            console.error('Error translating row:', error);
            // Restore original content for all loading cells
            loadingCells.forEach(({ translationText, originalContent }) => {
                translationText.innerHTML = originalContent;
            });
            showErrorToast('Translation failed', 'Translation Error');
        });
    }

    function showSourceLanguageModal(availableLanguages, callback) {
        // Get language display names
        const languageOptions = availableLanguages.map(langCode => {
            const langInfo = localizationData.languages.find(l => l.code === langCode);
            return {
                code: langCode,
                name: langInfo ? langInfo.name : langCode
            };
        }).sort((a, b) => a.name.localeCompare(b.name));

        const modalContent = `
            <div class="source-language-selection">
                <p>Select source language to translate from:</p>
                <div class="language-options">
                    ${languageOptions.map(lang => `
                        <button class="language-option-btn" data-language="${lang.code}">
                            ${lang.name}
                        </button>
                    `).join('')}
                </div>
            </div>
        `;

        showCustomDialog('Select Source Language', modalContent, (result) => {
            if (result) {
                callback(result);
            }
        });

        // Add event listeners to language option buttons
        setTimeout(() => {
            document.querySelectorAll('.language-option-btn').forEach(btn => {
                btn.addEventListener('click', function() {
                    const selectedLanguage = this.getAttribute('data-language');
                    // Close modal and call callback
                    const modal = document.querySelector('.modal-overlay');
                    if (modal) {
                        modal.remove();
                    }
                    callback(selectedLanguage);
                });
            });
        }, 100);
    }

    // Helper functions
    function getAvailableSourceLanguages(translations, excludeLanguage = null) {
        const availableLanguages = [];

        for (const [lang, text] of Object.entries(translations)) {
            if (lang !== excludeLanguage && text && text.trim() !== '') {
                availableLanguages.push(lang);
            }
        }

        return availableLanguages;
    }

    function hasTranslatableContent(translations) {
        return Object.values(translations).some(text => text && text.trim() !== '');
    }

    function getLanguageName(languageCode) {
        const langInfo = localizationData.languages.find(l => l.code === languageCode);
        return langInfo ? langInfo.name : languageCode;
    }

    function updateSingleTranslationFromAPI(key, language, newValue, allTranslations = null) {
        const formData = new FormData();
        formData.append('translation_key', key);

        if (allTranslations) {
            // Update all translations for the key
            formData.append('translations', JSON.stringify(allTranslations));
        } else {
            // Update single translation
            const keyData = localizationData.translationKeys.find(k => k.key === key);
            const existingTranslations = keyData ? { ...keyData.translations } : {};
            existingTranslations[language] = newValue;
            formData.append('translations', JSON.stringify(existingTranslations));
        }

        fetch('/api/localizations/update-translations', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (!data.success) {
                console.error('Failed to save translation:', data.error);
            }
        })
        .catch(error => {
            console.error('Error saving translation:', error);
        });
    }

    // Translation configuration functions
    function updateConfigurationUI() {
        const statusText = document.getElementById('config-status-text');
        const configBtn = document.getElementById('config-btn');

        if (translationConfig.enabled && translationConfig.hasApiKey) {
            statusText.textContent = 'Configured and ready';
            statusText.className = 'config-status-text configured';
            configBtn.innerHTML = '<i class="fas fa-edit"></i> Update';
        } else {
            statusText.textContent = 'Not configured';
            statusText.className = 'config-status-text not-configured';
            configBtn.innerHTML = '<i class="fas fa-cog"></i> Configure';
        }
    }

    function showTranslationConfigModal() {
        const isUpdate = translationConfig.enabled && translationConfig.hasApiKey;
        const modalTitle = isUpdate ? 'Update Translation Configuration' : 'Translation Configuration';
        const buttonText = isUpdate ? 'Update API Key' : 'Save API Key';
        const placeholderText = isUpdate ? 'Enter new Google Translate API key' : 'Enter your Google Translate API key';

        const modalContent = `
            <div class="translation-config-modal">
                <div class="form-group">
                    <label for="api-key-input">Google Translate API Key:</label>
                    <input type="password" id="api-key-input" placeholder="${placeholderText}" />
                    <div class="form-hint">
                        Your API key will be saved securely in the project's .fdawg-config file.
                        ${isUpdate ? 'Leave empty to keep current API key.' : ''}
                    </div>
                </div>
                <div class="config-help">
                    <h6>How to get Google Translate API Key:</h6>
                    <ol>
                        <li>Go to <a href="https://console.cloud.google.com/" target="_blank">Google Cloud Console</a></li>
                        <li>Create or select a project</li>
                        <li>Enable the "Cloud Translation API"</li>
                        <li>Go to "APIs & Services" → "Credentials"</li>
                        <li>Click "Create Credentials" → "API Key"</li>
                        <li>Copy the generated API key and paste it above</li>
                    </ol>
                </div>
                <div class="form-actions">
                    <button class="secondary-btn" onclick="closeConfigModal()">Cancel</button>
                    <button class="primary-btn" onclick="saveApiKey()">${buttonText}</button>
                </div>
            </div>
        `;

        showCustomDialog(modalTitle, modalContent);
    }

    function closeConfigModal() {
        const modal = document.querySelector('.modal-overlay');
        if (modal) {
            modal.remove();
        }
    }

    function saveApiKey() {
        const apiKeyInput = document.getElementById('api-key-input');
        const apiKey = apiKeyInput.value.trim();
        const isUpdate = translationConfig.enabled && translationConfig.hasApiKey;

        // For updates, allow empty API key (keeps current one)
        if (!isUpdate && !apiKey) {
            showErrorToast('Please enter a valid API key', 'Configuration Error');
            return;
        }

        // Show loading state
        const saveBtn = document.querySelector('.modal-overlay .primary-btn');
        const originalText = saveBtn.innerHTML;
        const loadingText = isUpdate ? '<i class="fas fa-spinner fa-spin"></i> Updating...' : '<i class="fas fa-spinner fa-spin"></i> Saving...';
        saveBtn.innerHTML = loadingText;
        saveBtn.disabled = true;

        const formData = new FormData();
        formData.append('api_key', apiKey);

        fetch('/api/localizations/update-api-key', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                // Update local config
                translationConfig.enabled = true;
                translationConfig.hasApiKey = true;

                // Update UI
                updateConfigurationUI();
                updateTranslateButtonStates();

                // Close modal
                closeConfigModal();

                const successMessage = isUpdate ? 'Google Translate API key updated successfully' : 'Google Translate API key saved successfully';
                showSuccessToast(successMessage, 'Configuration Updated');
            } else {
                showErrorToast(data.error || 'Failed to save API key', 'Configuration Error');
                saveBtn.innerHTML = originalText;
                saveBtn.disabled = false;
            }
        })
        .catch(error => {
            console.error('Error saving API key:', error);
            showErrorToast('Failed to save API key', 'Configuration Error');
            saveBtn.innerHTML = originalText;
            saveBtn.disabled = false;
        });
    }

    // Make functions globally available
    window.showTranslationConfigModal = showTranslationConfigModal;
    window.closeConfigModal = closeConfigModal;
    window.saveApiKey = saveApiKey;

});
