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
});
