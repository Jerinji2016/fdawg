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
        sidebarToggle.addEventListener('click', function() {
            sidebar.classList.toggle('expanded');
            sidebar.classList.toggle('collapsed');
            contentContainer.classList.toggle('expanded');
        });

        // Check screen size and set initial state
        function checkScreenSize() {
            if (window.innerWidth <= 768) {
                sidebar.classList.remove('collapsed');
                sidebar.classList.remove('expanded');
                contentContainer.classList.remove('expanded');
            } else if (window.innerWidth <= 992) {
                sidebar.classList.add('collapsed');
                sidebar.classList.remove('expanded');
                contentContainer.classList.add('expanded');
            } else {
                sidebar.classList.remove('collapsed');
                sidebar.classList.remove('expanded');
                contentContainer.classList.remove('expanded');
            }
        }

        // Run on load
        checkScreenSize();

        // Run on resize
        window.addEventListener('resize', checkScreenSize);
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
