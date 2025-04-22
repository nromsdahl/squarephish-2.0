(() => {
    'use strict'

    const THEME_STORAGE_KEY = 'theme'
    const themeToggler = document.getElementById('theme-toggler')
    const themeIcon = document.getElementById('theme-icon')
    const htmlElement = document.documentElement

    /**
     * Gets the preferred theme based on localStorage or OS setting.
     * @returns {'light' | 'dark'}
     */
    const getPreferredTheme = () => {
        const storedTheme = localStorage.getItem(THEME_STORAGE_KEY)
        if (storedTheme) {
            return storedTheme
        }
        // Check OS preference
        return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    }

    /**
     * Sets the theme on the <html> element, updates the icon, and saves preference.
     * @param {string} theme - The theme to set ('light' or 'dark').
     */
    const setTheme = (theme) => {
        if (theme !== 'light' && theme !== 'dark') {
            console.warn('Invalid theme specified:', theme);
            theme = 'light'; // Default to light if invalid
        }

        htmlElement.setAttribute('data-bs-theme', theme)
        updateIcon(theme)
        localStorage.setItem(THEME_STORAGE_KEY, theme)
    }

    /**
     * Updates the toggle button's icon based on the current theme.
     * @param {string} theme - The current theme ('light' or 'dark').
     */
    const updateIcon = (theme) => {
        if (theme === 'dark') {
            themeIcon.classList.remove('bi-moon-stars-fill');
            themeIcon.classList.add('bi-sun-fill');
            themeToggler.setAttribute('title', 'Switch to light mode');
        } else {
            themeIcon.classList.remove('bi-sun-fill');
            themeIcon.classList.add('bi-moon-stars-fill');
            themeToggler.setAttribute('title', 'Switch to dark mode');
        }
    }


    // --- Initialize Theme ---
    const currentTheme = getPreferredTheme()
    setTheme(currentTheme) // Apply theme immediately on load


    // --- Event Listener for Toggle Button ---
    themeToggler.addEventListener('click', () => {
        const currentTheme = htmlElement.getAttribute('data-bs-theme');
        const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
        setTheme(newTheme);
    });

    // --- Optional: Listen for OS theme changes ---
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', event => {
        // Only change if no theme is explicitly stored in localStorage
        if (!localStorage.getItem(THEME_STORAGE_KEY)) {
            const newColorScheme = event.matches ? "dark" : "light";
            setTheme(newColorScheme);
        }
    });
})() // Immediately Invoke Function Expression (IIFE) to scope variables