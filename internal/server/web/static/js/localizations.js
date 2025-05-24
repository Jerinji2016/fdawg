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

});
