/**
 * Toast Notification System
 *
 * This file contains functions for creating and managing toast notifications.
 */

// Toast counter for unique IDs
let toastCounter = 0;

/**
 * Show a toast notification
 * @param {string} message - The message to display
 * @param {string} type - The type of toast (success, error, info, warning)
 * @param {string} title - Optional title for the toast
 * @param {number} duration - Duration in milliseconds (default: 5000)
 * @returns {string} The ID of the created toast
 */
function showToast(message, type = 'info', title = '', duration = 5000) {
    const toastContainer = document.getElementById('toast-container');
    if (!toastContainer) return null;

    // Create a unique ID for this toast
    const toastId = `toast-${++toastCounter}`;

    // Create toast element
    const toast = document.createElement('div');
    toast.className = 'toast';
    toast.id = toastId;

    // Set icon based on type
    let iconClass = 'fa-info-circle';
    switch (type) {
        case 'success':
            iconClass = 'fa-check-circle';
            break;
        case 'error':
            iconClass = 'fa-exclamation-circle';
            break;
        case 'warning':
            iconClass = 'fa-exclamation-triangle';
            break;
    }

    // Create toast content
    toast.innerHTML = `
        <div class="toast-icon ${type}">
            <i class="fas ${iconClass}"></i>
        </div>
        <div class="toast-content">
            ${title ? `<div class="toast-title">${title}</div>` : ''}
            <div class="toast-message">${message}</div>
        </div>
        <button class="toast-close" aria-label="Close">
            <i class="fas fa-times"></i>
        </button>
        <div class="toast-progress ${type}"></div>
    `;

    // Add to container
    toastContainer.appendChild(toast);

    // Set up close button
    const closeButton = toast.querySelector('.toast-close');
    if (closeButton) {
        closeButton.addEventListener('click', () => {
            removeToast(toastId);
        });
    }

    // Auto-remove after duration
    if (duration > 0) {
        setTimeout(() => {
            removeToast(toastId);
        }, duration);
    }

    return toastId;
}

/**
 * Remove a toast notification with animation
 * @param {string} toastId - The ID of the toast to remove
 */
function removeToast(toastId) {
    const toast = document.getElementById(toastId);
    if (!toast) return;

    // Add exiting class for animation
    toast.classList.add('toast-exiting');

    // Remove after animation completes
    setTimeout(() => {
        if (toast.parentNode) {
            toast.parentNode.removeChild(toast);
        }
    }, 300); // Match the animation duration
}

/**
 * Show a success toast
 * @param {string} message - The message to display
 * @param {string} title - Optional title
 * @param {number} duration - Duration in milliseconds
 * @returns {string} The toast ID
 */
function showSuccessToast(message, title = 'Success', duration = 5000) {
    return showToast(message, 'success', title, duration);
}

/**
 * Show an error toast
 * @param {string} message - The message to display
 * @param {string} title - Optional title
 * @param {number} duration - Duration in milliseconds
 * @returns {string} The toast ID
 */
function showErrorToast(message, title = 'Error', duration = 5000) {
    return showToast(message, 'error', title, duration);
}

/**
 * Show an info toast
 * @param {string} message - The message to display
 * @param {string} title - Optional title
 * @param {number} duration - Duration in milliseconds
 * @returns {string} The toast ID
 */
function showInfoToast(message, title = 'Information', duration = 5000) {
    return showToast(message, 'info', title, duration);
}

/**
 * Show a warning toast
 * @param {string} message - The message to display
 * @param {string} title - Optional title
 * @param {number} duration - Duration in milliseconds
 * @returns {string} The toast ID
 */
function showWarningToast(message, title = 'Warning', duration = 5000) {
    return showToast(message, 'warning', title, duration);
}

/**
 * Show a confirmation toast with action buttons at the bottom
 * @param {string} message - The message to display
 * @param {string} title - The title for the toast
 * @param {Object} options - Configuration options
 * @param {string} options.confirmText - Text for the confirm button
 * @param {string} options.cancelText - Text for the cancel button
 * @param {string} options.confirmButtonClass - CSS class for the confirm button
 * @param {Function} options.onConfirm - Callback function when confirmed
 * @param {Function} options.onCancel - Callback function when canceled
 * @returns {string} The toast ID
 */
function showConfirmationToast(message, title, options = {}) {
    // Default options
    const defaultOptions = {
        confirmText: 'Confirm',
        cancelText: 'Cancel',
        confirmButtonClass: 'primary-btn',
        onConfirm: () => {},
        onCancel: () => {}
    };

    // Merge default options with provided options
    const config = { ...defaultOptions, ...options };

    // Create the toast
    const toastId = showToast(message, 'warning', title, 0); // Don't auto-dismiss

    // Get the toast element
    const toast = document.getElementById(toastId);
    if (!toast) return null;

    // Create action buttons container
    const actionsContainer = document.createElement('div');
    actionsContainer.className = 'toast-footer';

    // Cancel button
    const cancelBtn = document.createElement('button');
    cancelBtn.className = 'secondary-btn';
    cancelBtn.textContent = config.cancelText;
    cancelBtn.addEventListener('click', () => {
        removeToast(toastId);
        config.onCancel();
    });

    // Confirm button
    const confirmBtn = document.createElement('button');
    confirmBtn.className = config.confirmButtonClass;
    if (config.confirmText === 'Delete') {
        confirmBtn.style.backgroundColor = '#f44336';
    }
    confirmBtn.textContent = config.confirmText;
    confirmBtn.addEventListener('click', () => {
        removeToast(toastId);
        config.onConfirm();
    });

    // Add buttons to container
    actionsContainer.appendChild(cancelBtn);
    actionsContainer.appendChild(confirmBtn);

    // Add actions container to toast content
    const toastContent = toast.querySelector('.toast-content');
    if (toastContent) {
        toastContent.appendChild(actionsContainer);
    }

    // Remove progress bar
    const progressBar = toast.querySelector('.toast-progress');
    if (progressBar && progressBar.parentNode) {
        progressBar.parentNode.removeChild(progressBar);
    }

    return toastId;
}
