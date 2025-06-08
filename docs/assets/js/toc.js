/**
 * Table of Contents Generator and Scroll Highlighter
 * Automatically generates a right-side TOC and highlights active sections
 */

class TableOfContents {
    constructor() {
        this.tocContainer = null;
        this.headings = [];
        this.observer = null;
        this.activeLink = null;
        
        this.init();
    }

    init() {
        // Only initialize on pages with content
        if (!document.querySelector('.content-wrapper')) {
            return;
        }

        this.createTOCContainer();
        this.collectHeadings();
        
        if (this.headings.length > 0) {
            this.generateTOC();
            this.setupScrollObserver();
            this.setupClickHandlers();
        } else {
            // Hide TOC if no headings found
            this.tocContainer.style.display = 'none';
        }
    }

    createTOCContainer() {
        // Create TOC sidebar
        this.tocContainer = document.createElement('aside');
        this.tocContainer.className = 'toc-sidebar';
        this.tocContainer.innerHTML = `
            <div class="toc-header">On this page</div>
            <nav class="toc-nav">
                <ul id="toc-list"></ul>
            </nav>
        `;
        
        document.body.appendChild(this.tocContainer);
    }

    collectHeadings() {
        // Collect all headings from h2 to h6 (skip h1 as it's usually the page title)
        const headingSelectors = 'h2, h3, h4, h5, h6';
        const contentWrapper = document.querySelector('.content-wrapper');
        
        if (!contentWrapper) return;
        
        this.headings = Array.from(contentWrapper.querySelectorAll(headingSelectors))
            .filter(heading => {
                // Skip headings that are part of existing TOC or have no_toc class
                return !heading.classList.contains('no_toc') && 
                       !heading.closest('.no_toc') &&
                       heading.textContent.trim() !== 'Table of contents';
            })
            .map(heading => {
                // Ensure heading has an ID for linking
                if (!heading.id) {
                    heading.id = this.generateId(heading.textContent);
                }
                
                return {
                    element: heading,
                    id: heading.id,
                    text: heading.textContent.trim(),
                    level: parseInt(heading.tagName.charAt(1))
                };
            });
    }

    generateId(text) {
        // Generate a URL-friendly ID from heading text
        return text
            .toLowerCase()
            .replace(/[^\w\s-]/g, '') // Remove special characters
            .replace(/\s+/g, '-')     // Replace spaces with hyphens
            .replace(/-+/g, '-')      // Replace multiple hyphens with single
            .replace(/^-|-$/g, '');   // Remove leading/trailing hyphens
    }

    generateTOC() {
        const tocList = document.getElementById('toc-list');
        if (!tocList) return;

        tocList.innerHTML = '';

        this.headings.forEach(heading => {
            const listItem = document.createElement('li');
            listItem.className = `toc-h${heading.level}`;
            
            const link = document.createElement('a');
            link.href = `#${heading.id}`;
            link.textContent = heading.text;
            link.setAttribute('data-heading-id', heading.id);
            
            listItem.appendChild(link);
            tocList.appendChild(listItem);
        });
    }

    setupScrollObserver() {
        // Use Intersection Observer to detect which heading is currently visible
        const options = {
            root: null,
            rootMargin: '-60px 0px -80% 0px', // Account for header height and focus on top portion
            threshold: 0
        };

        this.observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                const headingId = entry.target.id;
                const tocLink = document.querySelector(`[data-heading-id="${headingId}"]`);
                
                if (entry.isIntersecting) {
                    // Remove active class from all links
                    document.querySelectorAll('.toc-nav a.active').forEach(link => {
                        link.classList.remove('active');
                    });
                    
                    // Add active class to current link
                    if (tocLink) {
                        tocLink.classList.add('active');
                        this.activeLink = tocLink;
                        
                        // Scroll TOC to show active item
                        this.scrollTOCToActive(tocLink);
                    }
                }
            });
        }, options);

        // Observe all headings
        this.headings.forEach(heading => {
            this.observer.observe(heading.element);
        });
    }

    scrollTOCToActive(activeLink) {
        if (!activeLink || !this.tocContainer) return;

        const tocNav = this.tocContainer.querySelector('.toc-nav');
        const linkRect = activeLink.getBoundingClientRect();
        const tocRect = tocNav.getBoundingClientRect();
        
        // Check if link is outside visible area
        if (linkRect.top < tocRect.top || linkRect.bottom > tocRect.bottom) {
            // Scroll to center the active link
            const scrollTop = activeLink.offsetTop - tocNav.offsetHeight / 2;
            tocNav.scrollTo({
                top: scrollTop,
                behavior: 'smooth'
            });
        }
    }

    setupClickHandlers() {
        // Handle TOC link clicks for smooth scrolling
        document.addEventListener('click', (e) => {
            if (e.target.matches('.toc-nav a[href^="#"]')) {
                e.preventDefault();
                
                const targetId = e.target.getAttribute('href').substring(1);
                const targetElement = document.getElementById(targetId);
                
                if (targetElement) {
                    // Smooth scroll to target with offset for header
                    const headerOffset = 80; // Account for sticky header + some padding
                    const elementPosition = targetElement.getBoundingClientRect().top;
                    const offsetPosition = elementPosition + window.pageYOffset - headerOffset;
                    
                    window.scrollTo({
                        top: offsetPosition,
                        behavior: 'smooth'
                    });
                    
                    // Update URL hash
                    history.pushState(null, null, `#${targetId}`);
                }
            }
        });
    }

    // Public method to refresh TOC (useful for dynamic content)
    refresh() {
        if (this.observer) {
            this.headings.forEach(heading => {
                this.observer.unobserve(heading.element);
            });
        }
        
        this.collectHeadings();
        
        if (this.headings.length > 0) {
            this.generateTOC();
            this.setupScrollObserver();
            this.tocContainer.style.display = 'block';
        } else {
            this.tocContainer.style.display = 'none';
        }
    }
}

// Initialize TOC when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    new TableOfContents();
});

// Handle page navigation (for single-page applications)
window.addEventListener('popstate', () => {
    // Small delay to ensure content is loaded
    setTimeout(() => {
        if (window.toc) {
            window.toc.refresh();
        }
    }, 100);
});

// Export for global access
window.TableOfContents = TableOfContents;
