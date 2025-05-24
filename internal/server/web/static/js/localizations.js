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
        const languages = localizationData.languages;

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
        const languages = localizationData.languages;
        const translationKeys = localizationData.translationKeys;

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
                    <td>${keyData.key}</td>
            `;

            // Add language columns
            languages.forEach(language => {
                const value = keyData.translations[language.code] || '';
                const isEmpty = !value.trim();
                rowsHTML += `
                    <td class="${isEmpty ? 'empty-translation' : ''}" title="${isEmpty ? 'Missing translation' : value}">
                        ${isEmpty ? '<em>Missing</em>' : value}
                    </td>
                `;
            });

            rowsHTML += `
                    <td>
                        <button class="table-btn edit-key-btn" data-key="${keyData.key}" title="Edit translations">
                            <i class="fas fa-edit"></i>
                        </button>
                        <button class="table-btn delete-key-btn" data-key="${keyData.key}" title="Delete key">
                            <i class="fas fa-trash"></i>
                        </button>
                    </td>
                </tr>
            `;
        });

        tbody.innerHTML = rowsHTML;

        // Add event listeners to table buttons
        addTableEventListeners();
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

    // Function to add event listeners to table buttons
    function addTableEventListeners() {
        // Edit buttons
        document.querySelectorAll('.edit-key-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                const key = this.getAttribute('data-key');
                showEditTranslationModal(key);
            });
        });

        // Delete buttons
        document.querySelectorAll('.delete-key-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                const key = this.getAttribute('data-key');
                showDeleteKeyConfirmation(key);
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

    // Modal and action functions will be added in the next part...

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

    function showEditTranslationModal(key) {
        const keyData = localizationData.translationKeys.find(k => k.key === key);
        if (!keyData) {
            showErrorToast('Translation key not found', 'Error');
            return;
        }

        // Create modal HTML for editing translations
        const languages = localizationData.languages;
        let languageInputsHTML = '';

        languages.forEach(language => {
            const value = keyData.translations[language.code] || '';
            languageInputsHTML += `
                <div class="form-group">
                    <label for="translation-${language.code}">${language.name} (${language.code}):</label>
                    <input type="text" id="translation-${language.code}" value="${value}" placeholder="Enter translation for ${language.name}">
                </div>
            `;
        });

        const modalHTML = `
            <div class="modal-overlay">
                <div class="modal-content">
                    <div class="modal-header">
                        <h3>Edit Translations for "${key}"</h3>
                        <button class="modal-close">&times;</button>
                    </div>
                    <div class="modal-body">
                        <form id="edit-translation-form">
                            ${languageInputsHTML}
                            <div class="form-actions">
                                <button type="button" class="secondary-btn cancel-btn">Cancel</button>
                                <button type="submit" class="primary-btn">Save Changes</button>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        `;

        // Add modal to the DOM
        document.body.insertAdjacentHTML('beforeend', modalHTML);

        // Get modal elements
        const modal = document.querySelector('.modal-overlay');
        const closeBtn = modal.querySelector('.modal-close');
        const cancelBtn = modal.querySelector('.cancel-btn');
        const form = modal.querySelector('#edit-translation-form');

        // Close modal function
        function closeModal() {
            modal.remove();
        }

        // Handle close events
        closeBtn.addEventListener('click', closeModal);
        cancelBtn.addEventListener('click', closeModal);

        // Handle form submission
        form.addEventListener('submit', function(e) {
            e.preventDefault();

            // Collect translation values
            const translations = {};
            languages.forEach(language => {
                const input = document.getElementById(`translation-${language.code}`);
                translations[language.code] = input.value.trim();
            });

            // Update translations
            updateTranslations(key, translations);
            closeModal();
        });

        // Close modal when clicking outside
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeModal();
            }
        });

        // Handle Escape key
        document.addEventListener('keydown', function escapeHandler(e) {
            if (e.key === 'Escape') {
                document.removeEventListener('keydown', escapeHandler);
                closeModal();
            }
        });
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

    function updateTranslations(key, translations) {
        const formData = new FormData();
        formData.append('translation_key', key);
        formData.append('translations', JSON.stringify(translations));

        fetch('/api/localizations/update-translations', {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showSuccessToast(`Translations for "${key}" updated successfully`, 'Success');
                loadLocalizationData(); // Reload data
            } else {
                showErrorToast(data.error || 'Failed to update translations', 'Error');
            }
        })
        .catch(error => {
            console.error('Error updating translations:', error);
            showErrorToast('Failed to update translations', 'Error');
        });
    }
});
