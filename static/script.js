document.addEventListener('DOMContentLoaded', () => {
    const addToCartForms = document.querySelectorAll('form[action="/cart/add"]');

    addToCartForms.forEach(form => {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const button = form.querySelector('.buy-btn');

            if (button) {
                button.style.transform = 'scale(0.95)';
                setTimeout(() => {
                    button.style.transform = 'scale(1)';
                }, 150);
            }

            const formData = new FormData(form);

            try {
                const response = await fetch('/cart/add', {
                    method: 'POST',
                    body: formData
                });

                if (response.ok) {
                    showDynamicToast('Товар успешно добавлен в корзину!', 'success');
                } else {
                    showDynamicToast('Не удалось добавить товар', 'error');
                }
            } catch (error) {
                console.error('Ошибка сети:', error);
                showDynamicToast('Ошибка соединения с сервером', 'error');
            }
        });
    });
});

function showDynamicToast(message, type = 'success') {
    const existingToast = document.querySelector('.cyber-toast');
    if (existingToast) {
        existingToast.remove();
    }

    const toast = document.createElement('div');
    toast.className = `cyber-toast ${type}`;

    const icon = type === 'success' ? '⚡' : '⚠️';
    const accentColor = type === 'success' ? '#10b981' : '#f43f5e';

    toast.innerHTML = `
        <span style="font-size: 18px;">${icon}</span>
        <div style="flex: 1;">
            <div style="font-weight: 600; color: #f8fafc; font-size: 13px;">${type === 'success' ? 'Успешно' : 'Ошибка'}</div>
            <div style="color: #94a3b8; font-size: 12px; margin-top: 2px;">${message}</div>
        </div>
    `;

    Object.assign(toast.style, {
        position: 'fixed',
        bottom: '30px',
        right: '30px',
        display: 'flex',
        alignItems: 'center',
        gap: '12px',
        backgroundColor: 'rgba(15, 23, 42, 0.85)',
        backdropFilter: 'blur(16px)',
        border: `1px solid ${accentColor}`,
        boxShadow: `0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 0 15px ${accentColor}40`,
        color: '#ffffff',
        padding: '16px 20px',
        borderRadius: '14px',
        zIndex: '10000',
        fontFamily: 'inherit',
        opacity: '0',
        transform: 'translateY(20px) scale(0.95)',
        transition: 'all 0.35s cubic-bezier(0.16, 1, 0.3, 1)'
    });

    document.body.appendChild(toast);

    requestAnimationFrame(() => {
        toast.style.opacity = '1';
        toast.style.transform = 'translateY(0) scale(1)';
    });

    setTimeout(() => {
        toast.style.opacity = '0';
        toast.style.transform = 'translateY(10px) scale(0.95)';
        setTimeout(() => toast.remove(), 350);
    }, 3200);
}
