// Global state
let currentUser = null;
let currentToken = null;
let items = [];
let auditHistory = [];

// API Configuration
const API_BASE = '/api';

// Initialize app
document.addEventListener('DOMContentLoaded', function() {
    checkAuthStatus();
});

// Authentication functions
async function login(event) {
    event.preventDefault();
    
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;
    
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            currentToken = data.token;
            const payload = parseJWT(data.token);
            currentUser = {
                id: payload.user_id,
                username: payload.user_name,
                role: payload.role
            };
            
            localStorage.setItem('token', data.token);
            showDashboard();
            showToast('Успешный вход в систему!', 'success');
        } else {
            showToast(data.error || 'Ошибка входа', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Login error:', error);
    } finally {
        showLoading(false);
    }
}

async function register(event) {
    event.preventDefault();
    
    const username = document.getElementById('regUsername').value;
    const password = document.getElementById('regPassword').value;
    const role = document.getElementById('regRole').value;
    
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ username, password, role })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showToast('Регистрация успешна! Теперь вы можете войти в систему.', 'success');
            showLogin();
            // Clear form
            document.getElementById('regUsername').value = '';
            document.getElementById('regPassword').value = '';
            document.getElementById('regRole').value = '';
        } else {
            showToast(data.error || 'Ошибка регистрации', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Register error:', error);
    } finally {
        showLoading(false);
    }
}

function logout() {
    currentUser = null;
    currentToken = null;
    localStorage.removeItem('token');
    showAuth();
    showToast('Вы вышли из системы', 'success');
}

function checkAuthStatus() {
    const token = localStorage.getItem('token');
    if (token) {
        try {
            const payload = parseJWT(token);
            // Check if token is not expired
            if (payload.exp * 1000 > Date.now()) {
                currentToken = token;
                currentUser = {
                    id: payload.user_id,
                    username: payload.user_name,
                    role: payload.role
                };
                showDashboard();
                return;
            }
        } catch (error) {
            console.error('Invalid token:', error);
        }
        localStorage.removeItem('token');
    }
    showAuth();
}

function parseJWT(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
}

// UI Navigation
function showAuth() {
    document.getElementById('authSection').style.display = 'block';
    document.getElementById('dashboardSection').style.display = 'none';
    document.getElementById('userInfo').style.display = 'none';
}

function showDashboard() {
    document.getElementById('authSection').style.display = 'none';
    document.getElementById('dashboardSection').style.display = 'block';
    document.getElementById('userInfo').style.display = 'flex';
    
    // Update user info
    document.getElementById('userRole').textContent = getRoleDisplayName(currentUser.role);
    document.getElementById('userName').textContent = currentUser.username;
    
    // Load initial data
    loadItems();
    loadAuditHistory();
}

function showLogin() {
    document.getElementById('loginForm').style.display = 'block';
    document.getElementById('registerForm').style.display = 'none';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelector('.tab-btn').classList.add('active');
}

function showRegister() {
    document.getElementById('loginForm').style.display = 'none';
    document.getElementById('registerForm').style.display = 'block';
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    document.querySelectorAll('.tab-btn')[1].classList.add('active');
}

function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.style.display = 'none';
    });
    document.querySelectorAll('.nav-tab').forEach(tab => {
        tab.classList.remove('active');
    });
    
    // Show selected tab
    document.getElementById(tabName + 'Tab').style.display = 'block';
    event.target.classList.add('active');
    
    // Load data for the tab
    if (tabName === 'items') {
        loadItems();
    } else if (tabName === 'audit') {
        loadAuditHistory();
    }
}

// Items Management
async function loadItems() {
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE}/items`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
        
        if (response.ok) {
            items = await response.json();
            renderItems();
        } else {
            showToast('Ошибка загрузки товаров', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Load items error:', error);
    } finally {
        showLoading(false);
    }
}

function renderItems() {
    const container = document.getElementById('itemsGrid');
    container.innerHTML = '';
    
    if (items.length === 0) {
        container.innerHTML = '<div class="text-center" style="grid-column: 1 / -1; padding: 2rem; color: #6c757d;">Товары не найдены</div>';
        return;
    }
    
    items.forEach(item => {
        const itemCard = document.createElement('div');
        itemCard.className = 'item-card';
        itemCard.innerHTML = `
            <div class="item-header">
                <div>
                    <div class="item-name">${escapeHtml(item.name)}</div>
                    <div class="item-id">ID: ${item.id}</div>
                </div>
            </div>
            <div class="item-description">${escapeHtml(item.description || 'Нет описания')}</div>
            <div class="item-quantity">Количество: ${item.quantity}</div>
            <div class="item-actions">
                ${canEdit() ? `<button class="btn btn-secondary" onclick="editItem(${item.id})">
                    <i class="fas fa-edit"></i>
                    Редактировать
                </button>` : ''}
                ${canDelete() ? `<button class="btn btn-danger" onclick="deleteItem(${item.id})">
                    <i class="fas fa-trash"></i>
                    Удалить
                </button>` : ''}
            </div>
        `;
        container.appendChild(itemCard);
    });
}

function showAddItemModal() {
    if (!canEdit()) {
        showToast('Недостаточно прав для добавления товаров', 'error');
        return;
    }
    
    document.getElementById('modalTitle').textContent = 'Добавить товар';
    document.getElementById('itemForm').reset();
    document.getElementById('itemModal').classList.add('show');
}

function editItem(itemId) {
    if (!canEdit()) {
        showToast('Недостаточно прав для редактирования товаров', 'error');
        return;
    }
    
    const item = items.find(i => i.id === itemId);
    if (!item) return;
    
    document.getElementById('modalTitle').textContent = 'Редактировать товар';
    document.getElementById('itemName').value = item.name;
    document.getElementById('itemDescription').value = item.description || '';
    document.getElementById('itemQuantity').value = item.quantity;
    document.getElementById('itemForm').dataset.itemId = itemId;
    document.getElementById('itemModal').classList.add('show');
}

async function saveItem(event) {
    event.preventDefault();
    
    const formData = {
        name: document.getElementById('itemName').value,
        description: document.getElementById('itemDescription').value,
        quantity: parseInt(document.getElementById('itemQuantity').value)
    };
    
    const itemId = document.getElementById('itemForm').dataset.itemId;
    
    try {
        showLoading(true);
        
        let response;
        if (itemId) {
            // Update existing item
            response = await fetch(`${API_BASE}/items/${itemId}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${currentToken}`
                },
                body: JSON.stringify(formData)
            });
        } else {
            // Create new item
            response = await fetch(`${API_BASE}/items`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${currentToken}`
                },
                body: JSON.stringify(formData)
            });
        }
        
        if (response.ok) {
            showToast(itemId ? 'Товар обновлен' : 'Товар добавлен', 'success');
            closeItemModal();
            loadItems();
            loadAuditHistory(); // Refresh audit history
        } else {
            const data = await response.json();
            showToast(data.error || 'Ошибка сохранения', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Save item error:', error);
    } finally {
        showLoading(false);
    }
}

async function deleteItem(itemId) {
    if (!canDelete()) {
        showToast('Недостаточно прав для удаления товаров', 'error');
        return;
    }
    
    if (!confirm('Вы уверены, что хотите удалить этот товар?')) {
        return;
    }
    
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE}/items/${itemId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
        
        if (response.ok) {
            showToast('Товар удален', 'success');
            loadItems();
            loadAuditHistory(); // Refresh audit history
        } else {
            const data = await response.json();
            showToast(data.error || 'Ошибка удаления', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Delete item error:', error);
    } finally {
        showLoading(false);
    }
}

function closeItemModal() {
    document.getElementById('itemModal').classList.remove('show');
    document.getElementById('itemForm').reset();
    delete document.getElementById('itemForm').dataset.itemId;
}

// Audit History
async function loadAuditHistory() {
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE}/audit`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
        
        if (response.ok) {
            auditHistory = await response.json();
            renderAuditHistory();
            updateItemFilter();
        } else {
            showToast('Ошибка загрузки истории', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Load audit history error:', error);
    } finally {
        showLoading(false);
    }
}

function renderAuditHistory() {
    const tbody = document.getElementById('auditTableBody');
    tbody.innerHTML = '';
    
    if (auditHistory.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center">История изменений пуста</td></tr>';
        return;
    }
    
    auditHistory.forEach(record => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${record.id}</td>
            <td>${record.item_id === 0 ? '-' : record.item_id}</td>
            <td><span class="action-badge action-${record.action.toLowerCase()}">${record.action}</span></td>
            <td>${record.user_id}</td>
            <td>${formatDate(record.created_at)}</td>
            <td>${record.old_value ? formatJSON(record.old_value) : '-'}</td>
            <td>${record.new_value ? formatJSON(record.new_value) : '-'}</td>
        `;
        tbody.appendChild(row);
    });
}

function updateItemFilter() {
    const filter = document.getElementById('itemFilter');
    const uniqueItemIds = [...new Set(auditHistory.map(record => record.item_id))].filter(id => id && id !== 0);
    
    filter.innerHTML = '<option value="">Все товары</option>';
    uniqueItemIds.forEach(itemId => {
        const option = document.createElement('option');
        option.value = itemId;
        option.textContent = `Товар ID: ${itemId}`;
        filter.appendChild(option);
    });
}

function filterAuditHistory() {
    const itemId = document.getElementById('itemFilter').value;
    const filteredHistory = itemId ? 
        auditHistory.filter(record => record.item_id == itemId) : 
        auditHistory;
    
    const tbody = document.getElementById('auditTableBody');
    tbody.innerHTML = '';
    
    if (filteredHistory.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center">Нет записей для выбранного фильтра</td></tr>';
        return;
    }
    
    filteredHistory.forEach(record => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${record.id}</td>
            <td>${record.item_id}</td>
            <td><span class="action-badge action-${record.action.toLowerCase()}">${record.action}</span></td>
            <td>${record.user_id}</td>
            <td>${formatDate(record.created_at)}</td>
            <td>${record.old_value ? formatJSON(record.old_value) : '-'}</td>
            <td>${record.new_value ? formatJSON(record.new_value) : '-'}</td>
        `;
        tbody.appendChild(row);
    });
}

async function exportHistory() {
    try {
        const response = await fetch(`${API_BASE}/audit/export`, {
            headers: {
                'Authorization': `Bearer ${currentToken}`
            }
        });
        
        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `audit_history_${new Date().toISOString().split('T')[0]}.csv`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
            showToast('История экспортирована', 'success');
        } else {
            showToast('Ошибка экспорта', 'error');
        }
    } catch (error) {
        showToast('Ошибка соединения', 'error');
        console.error('Export error:', error);
    }
}

// Permission checks
function canEdit() {
    return currentUser && ['admin', 'manager'].includes(currentUser.role);
}

function canDelete() {
    return currentUser && currentUser.role === 'admin';
}

function canViewAudit() {
    return currentUser && ['admin', 'manager', 'viewer'].includes(currentUser.role);
}

// Utility functions
function showLoading(show) {
    document.getElementById('loadingSpinner').style.display = show ? 'flex' : 'none';
}

function showToast(message, type = 'info') {
    const container = document.getElementById('toastContainer');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 5000);
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString('ru-RU');
}

function formatJSON(jsonString) {
    try {
        const obj = typeof jsonString === 'string' ? JSON.parse(jsonString) : jsonString;
        return JSON.stringify(obj, null, 2);
    } catch (error) {
        return jsonString;
    }
}

function getRoleDisplayName(role) {
    const roleNames = {
        'admin': 'Администратор',
        'manager': 'Менеджер',
        'viewer': 'Наблюдатель'
    };
    return roleNames[role] || role;
}

// Close modal when clicking outside
document.addEventListener('click', function(event) {
    const modal = document.getElementById('itemModal');
    if (event.target === modal) {
        closeItemModal();
    }
});

// Close modal with Escape key
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        closeItemModal();
    }
});
