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
        setTimeout(() => {
            card.style.opacity = '1';
            card.style.transform = 'translateY(0)';
        }, 50 * index);
    });
});
