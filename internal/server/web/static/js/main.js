// Main JavaScript file for the web interface
document.addEventListener('DOMContentLoaded', function() {
    console.log("Flutter Project Manager loaded");

    // Format dependency versions for better display
    const depVersions = document.querySelectorAll('.dep-version');

    depVersions.forEach(function(element) {
        const content = element.textContent;
        if (content.includes('map') || content.includes('object')) {
            element.textContent = '(complex dependency)';
        }
    });

    // Add animation classes to cards for staggered entrance
    const cards = document.querySelectorAll('.info-card, .dependency-item, .asset-list li');
    cards.forEach((card, index) => {
        card.style.opacity = '0';
        card.style.transform = 'translateY(20px)';
        setTimeout(() => {
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, 50 * index);
    });

    // Add click event for section headers to toggle visibility
    const sectionHeaders = document.querySelectorAll('.info-section h3');
    sectionHeaders.forEach(header => {
        header.style.cursor = 'pointer';
        header.addEventListener('click', function() {
            const content = this.nextElementSibling;
            if (content.style.display === 'none') {
                content.style.display = '';
                this.querySelector('i.fas.fa-ellipsis-h').classList.remove('fa-ellipsis-h');
                this.querySelector('i:last-child').classList.add('fa-ellipsis-h');
            } else {
                content.style.display = 'none';
                this.querySelector('i.fas.fa-ellipsis-h').classList.remove('fa-ellipsis-h');
                this.querySelector('i:last-child').classList.add('fa-ellipsis-v');
            }
        });
    });

    // Sidebar toggle functionality
    const sidebarToggle = document.getElementById('sidebar-toggle');
    const sidebar = document.getElementById('sidebar');
    const contentContainer = document.querySelector('.content-container');

    if (sidebarToggle && sidebar) {
        // Function to save sidebar state to localStorage
        function saveSidebarState(isCollapsed) {
            localStorage.setItem('sidebarCollapsed', isCollapsed);
        }

        // Function to load sidebar state from localStorage
        function loadSidebarState() {
            return localStorage.getItem('sidebarCollapsed') === 'true';
        }

        // Function to apply sidebar state
        function applySidebarState(isCollapsed) {
            if (isCollapsed) {
                sidebar.classList.add('collapsed');
                contentContainer.classList.add('expanded');
            } else {
                sidebar.classList.remove('collapsed');
                contentContainer.classList.remove('expanded');
            }
            sidebar.classList.remove('expanded');
        }

        // Toggle sidebar on button click
        sidebarToggle.addEventListener('click', function() {
            const isCurrentlyCollapsed = sidebar.classList.contains('collapsed');
            const willBeCollapsed = !isCurrentlyCollapsed;

            sidebar.classList.toggle('expanded');
            sidebar.classList.toggle('collapsed');
            contentContainer.classList.toggle('expanded');

            // Save the new state
            saveSidebarState(willBeCollapsed);
        });

        // Check screen size and set initial state
        function checkScreenSize() {
            const savedState = loadSidebarState();

            if (window.innerWidth <= 768) {
                // On mobile, always start with sidebar hidden
                sidebar.classList.remove('collapsed');
                sidebar.classList.remove('expanded');
                contentContainer.classList.remove('expanded');

                // Add a click handler to close sidebar when clicking outside on mobile
                if (!window.sidebarClickOutsideHandler) {
                    window.sidebarClickOutsideHandler = true;
                    document.addEventListener('click', function(e) {
                        // Only apply this on mobile
                        if (window.innerWidth <= 768) {
                            // If sidebar is expanded and click is outside sidebar
                            if (sidebar.classList.contains('expanded') &&
                                !sidebar.contains(e.target) &&
                                e.target !== sidebarToggle &&
                                !sidebarToggle.contains(e.target)) {
                                sidebar.classList.remove('expanded');
                            }
                        }
                    });
                }
            } else if (window.innerWidth <= 992) {
                // On tablet, use saved state or default to collapsed
                if (savedState !== null) {
                    applySidebarState(savedState);
                } else {
                    sidebar.classList.add('collapsed');
                    contentContainer.classList.add('expanded');
                }
            } else {
                // On desktop, use saved state or default to expanded
                if (savedState !== null) {
                    applySidebarState(savedState);
                } else {
                    sidebar.classList.remove('collapsed');
                    contentContainer.classList.remove('expanded');
                }
            }
        }

        // Run on load
        checkScreenSize();

        // Run on resize
        window.addEventListener('resize', checkScreenSize);

        // Enhance tooltip behavior for better UX
        const navItems = document.querySelectorAll('.nav-item a');
        navItems.forEach(item => {
            // Add a small delay before showing tooltip to prevent accidental triggers
            let tooltipTimer;

            item.addEventListener('mouseenter', function() {
                if (sidebar.classList.contains('collapsed')) {
                    // Add a small delay before showing tooltip
                    tooltipTimer = setTimeout(() => {
                        this.classList.add('tooltip-active');
                    }, 200);
                }
            });

            item.addEventListener('mouseleave', function() {
                // Clear the timer if mouse leaves before tooltip is shown
                clearTimeout(tooltipTimer);
                this.classList.remove('tooltip-active');
            });

            // Modify navigation links to preserve sidebar state
            item.addEventListener('click', function(e) {
                // Only intercept if it's a navigation link (not an external link)
                if (this.getAttribute('href').startsWith('/') || this.getAttribute('href') === '#') {
                    e.preventDefault();

                    // On mobile, close the sidebar when a link is clicked
                    if (window.innerWidth <= 768) {
                        sidebar.classList.remove('expanded');
                    }

                    // Save current sidebar state before navigation
                    const isCollapsed = sidebar.classList.contains('collapsed');
                    saveSidebarState(isCollapsed);

                    // Navigate to the link
                    window.location.href = this.getAttribute('href');
                }
            });
        });
    }

    // Tab functionality
    const tabButtons = document.querySelectorAll('.tab-btn');

    tabButtons.forEach(button => {
        button.addEventListener('click', function() {
            const tabId = this.getAttribute('data-tab');

            // Remove active class from all buttons and panes
            document.querySelectorAll('.tab-btn').forEach(btn => {
                btn.classList.remove('active');
            });

            document.querySelectorAll('.tab-pane').forEach(pane => {
                pane.classList.remove('active');
            });

            // Add active class to clicked button and corresponding pane
            this.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        });
    });
});
