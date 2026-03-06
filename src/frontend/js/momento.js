function showCustomTab(tabId) {
    document.querySelectorAll('.custom-tab-content').forEach(tab => tab.style.display = 'none');
    document.querySelectorAll('.custom-nav-item').forEach(item => item.classList.remove('active'));
    document.getElementById(tabId).style.display = 'block';
    document.querySelector(`[onclick="showCustomTab('${tabId}')"]`).classList.add('active');
  }
