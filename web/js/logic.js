(function() {
    'use strict';

    // Configuration
    const API_BASE = window.location.origin + '/api';
    // const API_BASE = 'http://localhost:31979/api';

    // DOM Elements
    const taggedKvTemplate = document.getElementById("tagged-key-value-template");
    const taggedKvContainer = document.getElementById("tagged-key-values-container");
    const searchInput = document.getElementById("search-input");
    const newKeyInput = document.getElementById("new-key-input");
    const newValueInput = document.getElementById("new-value-input");
    const createItemBtn = document.getElementById("create-item-btn");

    // State
    let lastSearchTags = [];
    let searchDebounceTimer = null;

    // Utility: Debounce function
    function debounce(func, delay) {
        return function(...args) {
            clearTimeout(searchDebounceTimer);
            searchDebounceTimer = setTimeout(() => func.apply(this, args), delay);
        };
    }

    // Utility: Escape HTML to prevent XSS
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Utility: Handle API response
    async function handleResponse(response) {
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || `HTTP error! status: ${response.status}`);
        }
        return response;
    }

    // Utility: Show notification
    function showNotification(message, type = 'info') {
        // Using alert for now, but this should be replaced with a toast notification
        if (type === 'error') {
            alert(`Error: ${message}`);
        } else {
            console.log(message);
        }
    }

    // API: Create a new item
    async function createItem(key, value) {
        const response = await fetch(`${API_BASE}/keys`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ key, value })
        });
        return handleResponse(response);
    }

    // API: Add a tag to an item
    async function addTag(key, tag) {
        const response = await fetch(`${API_BASE}/tags`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ key, tag })
        });
        return handleResponse(response);
    }

    // API: Remove a tag from an item
    async function removeTag(key, tag) {
        const response = await fetch(`${API_BASE}/tags/${encodeURIComponent(tag)}/${encodeURIComponent(key)}`, {
            method: 'DELETE'
        });
        return handleResponse(response);
    }

    // API: Delete an item
    async function deleteItem(key) {
        const response = await fetch(`${API_BASE}/keys/${encodeURIComponent(key)}`, {
            method: 'DELETE'
        });
        return handleResponse(response);
    }

    // API: Search for items by tags
    async function searchByTags(tags) {
        const queryString = tags.length > 0 ? `?tags=${tags.join(',')}` : '';
        const response = await fetch(`${API_BASE}/keys${queryString}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        await handleResponse(response);
        return response.json();
    }

    // Create a tag badge with remove button
    function createTagBadge(key, tag) {
        const tagBadge = document.createElement('span');
        tagBadge.className = 'tag-badge';

        const tagText = document.createElement('span');
        tagText.textContent = tag;

        const removeBtn = document.createElement('button');
        removeBtn.className = 'tag-remove-btn';
        removeBtn.textContent = '×';
        removeBtn.type = 'button';
        removeBtn.setAttribute('aria-label', `Remove tag ${tag}`);

        removeBtn.addEventListener('click', async () => {
            try {
                await removeTag(key, tag);
                tagBadge.remove();
                console.log(`Removed tag "${tag}" from key "${key}"`);
            } catch (error) {
                console.error('Error removing tag:', error);
                showNotification(`Failed to remove tag: ${error.message}`, 'error');
            }
        });

        tagBadge.appendChild(tagText);
        tagBadge.appendChild(removeBtn);

        return tagBadge;
    }

    // Render a tagged key-value item
    function renderItem(item) {
        const clone = taggedKvTemplate.content.cloneNode(true);

        // Set key and value using textContent to prevent XSS
        const titleElement = clone.querySelector(".tkv-title");
        const valueElement = clone.querySelector(".tkv-value");
        titleElement.textContent = item.key;
        valueElement.textContent = item.value;

        // Setup delete button
        const deleteBtn = clone.querySelector(".tkv-delete-btn");
        deleteBtn.addEventListener('click', async () => {
            if (!confirm(`Are you sure you want to delete "${item.key}"?`)) {
                return;
            }

            try {
                await deleteItem(item.key);
                // Find and remove the item from DOM
                const tkvElement = deleteBtn.closest('.tkv');
                if (tkvElement) {
                    tkvElement.remove();
                }
                console.log(`Deleted item "${item.key}"`);
            } catch (error) {
                console.error('Error deleting item:', error);
                showNotification(`Failed to delete item: ${error.message}`, 'error');
            }
        });

        // Render existing tags
        const tagsContainer = clone.querySelector(".tkv-tags");
        const existingTags = item.tags || [];
        existingTags.forEach((tag) => {
            tagsContainer.appendChild(createTagBadge(item.key, tag));
        });

        // Setup add tag functionality
        const tagInput = clone.querySelector(".tkv-tag-input");
        const addTagBtn = clone.querySelector(".tkv-add-tag-btn");

        const handleAddTag = async () => {
            const newTag = tagInput.value.trim();
            if (!newTag) {
                showNotification('Please enter a tag', 'error');
                return;
            }

            // Check for duplicate tags
            if (existingTags.includes(newTag)) {
                showNotification('Tag already exists', 'error');
                return;
            }

            try {
                await addTag(item.key, newTag);
                tagsContainer.appendChild(createTagBadge(item.key, newTag));
                existingTags.push(newTag);
                tagInput.value = '';
                console.log(`Added tag "${newTag}" to key "${item.key}"`);
            } catch (error) {
                console.error('Error adding tag:', error);
                showNotification(`Failed to add tag: ${error.message}`, 'error');
            }
        };

        addTagBtn.addEventListener('click', handleAddTag);

        // Allow Enter key to add tag
        tagInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                handleAddTag();
            }
        });

        return clone;
    }

    // Clear and render search results
    function renderSearchResults(items) {
        // Clear existing entries efficiently
        taggedKvContainer.innerHTML = '';

        // Render items
        items.forEach((item) => {
            taggedKvContainer.appendChild(renderItem(item));
        });
    }

    // Handle search
    async function handleSearch() {
        const searchValue = searchInput.value.trim();
        const tags = searchValue === "*" || searchValue === "" ? [] : searchValue.split(" ").filter(t => t);

        // Avoid duplicate searches
        if (tags.toString() === lastSearchTags.toString()) {
            return;
        }

        console.log("Searching for tags:", tags);
        lastSearchTags = tags;

        try {
            const data = await searchByTags(tags);
            console.log("Received data:", data);
            renderSearchResults(data);
        } catch (error) {
            console.error('Search error:', error);
            showNotification(`Search failed: ${error.message}`, 'error');
            renderSearchResults([]);
        }
    }

    // Handle create item
    async function handleCreateItem() {
        const key = newKeyInput.value.trim();
        const value = newValueInput.value.trim();

        if (!key) {
            showNotification('Please enter a key', 'error');
            return;
        }

        if (!value) {
            showNotification('Please enter a value', 'error');
            return;
        }

        try {
            await createItem(key, value);

            // Clear inputs
            newKeyInput.value = '';
            newValueInput.value = '';

            // Add the new item to the display
            const newItem = {
                key: key,
                value: value,
                tags: []
            };

            taggedKvContainer.insertBefore(
                renderItem(newItem),
                taggedKvContainer.firstChild
            );

            console.log(`Created item with key "${key}"`);
            showNotification('Item created successfully!');
        } catch (error) {
            console.error('Error creating item:', error);
            showNotification(`Failed to create item: ${error.message}`, 'error');
        }
    }

    // Event Listeners
    // Debounced search on input change (400ms delay)
    const debouncedSearch = debounce(handleSearch, 400);

    searchInput.addEventListener("input", () => {
        debouncedSearch();
    });

    // Allow Enter to trigger immediate search
    searchInput.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            event.preventDefault();
            clearTimeout(searchDebounceTimer);
            handleSearch();
        }
    });

    createItemBtn.addEventListener('click', handleCreateItem);

    // Allow Enter key in value input to create item
    newValueInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleCreateItem();
        }
    });

    // Allow Enter key in key input to focus value input
    newKeyInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            newValueInput.focus();
        }
    });

    // Initial load - show all items
    handleSearch();

})();
