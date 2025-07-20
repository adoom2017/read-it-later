import React, { useState, useEffect } from 'react';
import axios from 'axios';
import ArticleList from './components/ArticleList';
import SearchBar from './components/SearchBar';
import AuthModal from './components/AuthModal';
import { useAuth, api } from './hooks/useAuth';
import './App.css';
import './components/AuthModal.css';

function App() {
  const [url, setUrl] = useState('');
  const [articles, setArticles] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState(null);
  const [isSearching, setIsSearching] = useState(false);
  const [searchResults, setSearchResults] = useState([]);
  const [currentSearchQuery, setCurrentSearchQuery] = useState('');
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [authMessage, setAuthMessage] = useState('');

  // 使用认证Hook
  const { user, loading: authLoading, error: authError, isAuthenticated, login, register, logout } = useAuth();

  // Fetch articles on component mount
  useEffect(() => {
    const fetchArticles = async () => {
      if (!isAuthenticated) return;
      
      try {
        const data = await api.getArticles();
        setArticles(data || []);
      } catch (err) {
        setError('Failed to fetch articles.');
        console.error(err);
      }
    };
    fetchArticles();
  }, [isAuthenticated]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!url) return;

    setLoading(true);
    setError('');

    try {
      const newArticle = await api.createArticle({ url });
      setArticles([newArticle, ...articles]);
      setUrl(''); // Clear input field
    } catch (err) {
      setError('Failed to add article. Please check the URL and try again.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // 删除文章功能
  const handleDeleteArticle = async (articleId) => {
    if (!window.confirm('Are you sure you want to delete this article?')) {
      return;
    }

    try {
      await api.deleteArticle(articleId);
      setArticles(articles.filter(article => article.id !== articleId));
      setError('');
    } catch (err) {
      setError('Failed to delete article.');
      console.error(err);
    }
  };

  // 添加标签功能
  const handleAddTag = async (articleId, tagName) => {
    if (!tagName.trim()) return;

    try {
      await api.addTagToArticle(articleId, tagName.trim());
      // 重新获取文章列表以显示更新的标签
      const data = await api.getArticles();
      setArticles(data || []);
      setError('');
      
      // 如果正在搜索，也更新搜索结果
      if (isSearching) {
        // 重新执行搜索以更新结果
        const searchType = currentSearchQuery.includes('#') ? 'tag' : 'title';
        await handleSearch(currentSearchQuery, searchType);
      }
    } catch (err) {
      setError('Failed to add tag.');
      console.error(err);
    }
  };

  // 删除标签功能
  const handleRemoveTag = async (articleId, tagId) => {
    if (!window.confirm('确定要删除这个标签吗？')) {
      return;
    }

    try {
      await api.removeTagFromArticle(articleId, tagId);
      // 重新获取文章列表以显示更新的标签
      const data = await api.getArticles();
      setArticles(data || []);
      setError('');
      
      // 如果正在搜索，也更新搜索结果
      if (isSearching) {
        const searchType = currentSearchQuery.includes('#') ? 'tag' : 'title';
        await handleSearch(currentSearchQuery, searchType);
      }
    } catch (err) {
      setError('Failed to remove tag.');
      console.error(err);
    }
  };

  // 查看文章详情功能
  const handleViewArticle = async (articleId) => {
    try {
      const data = await api.getArticle(articleId);
      setSelectedArticle(data);
      setError('');
    } catch (err) {
      setError('Failed to fetch article details.');
      console.error(err);
    }
  };

  const closeArticleDetail = () => {
    setSelectedArticle(null);
  };

  // 搜索功能
  const handleSearch = async (query, searchType) => {
    setLoading(true);
    setError('');
    
    try {
      const results = await api.searchArticles(query, searchType);
      setSearchResults(results || []);
      setIsSearching(true);
      setCurrentSearchQuery(query);
    } catch (err) {
      setError('搜索失败，请重试。');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  // 清除搜索
  const handleClearSearch = () => {
    setIsSearching(false);
    setSearchResults([]);
    setCurrentSearchQuery('');
    setError('');
  };

  // 认证相关处理函数
  const handleLogin = async (username, password) => {
    const result = await login(username, password);
    if (result.success) {
      setShowAuthModal(false);
      setAuthMessage('登录成功！');
      setTimeout(() => setAuthMessage(''), 3000);
    } else {
      setError(result.error);
    }
  };

  const handleRegister = async (username, email, password) => {
    const result = await register(username, email, password);
    if (result.success) {
      setShowAuthModal(false);
      setAuthMessage('注册成功！');
      setTimeout(() => setAuthMessage(''), 3000);
    } else {
      setError(result.error);
    }
  };

  const handleLogout = () => {
    logout();
    setArticles([]);
    setAuthMessage('已退出登录');
    setTimeout(() => setAuthMessage(''), 3000);
  };

  // 如果还在加载认证状态，显示loading
  if (authLoading) {
    return (
      <div className="App">
        <div className="loading-container">
          <p>加载中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="App">
      <header className="App-header">
        <h1>Read It Later</h1>
        <p>Save articles to read and summarize them.</p>
        
        {/* 用户认证状态显示 - 只在已登录时显示 */}
        {isAuthenticated && (
          <div className="auth-section">
            <div className="user-info">
              <div className="user-details">
                <p className="username">欢迎, {user?.username}</p>
                {user?.email && <p className="user-email">{user.email}</p>}
              </div>
              <button onClick={handleLogout} className="logout-btn">
                退出登录
              </button>
            </div>
          </div>
        )}
        
        {authMessage && <p className="success-message">{authMessage}</p>}
      </header>
      
      {/* 只有登录后才显示主要功能 */}
      {isAuthenticated ? (
        <main>
          <form onSubmit={handleSubmit} className="url-form">
            <input
              type="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              placeholder="https://example.com/article"
              required
              disabled={loading}
            />
            <button type="submit" disabled={loading}>
              {loading ? 'Saving...' : 'Save Article'}
            </button>
          </form>
          {error && <p className="error-message">{error}</p>}
          
          {/* 搜索组件 */}
          <SearchBar onSearch={handleSearch} onClear={handleClearSearch} />
        
          {/* 显示搜索结果或所有文章 */}
          {isSearching && (
            <div className="search-results-header">
              <h3>搜索结果: "{currentSearchQuery}" ({searchResults.length} 条结果)</h3>
              <button onClick={handleClearSearch} className="back-to-all-btn">
                返回所有文章
              </button>
            </div>
          )}
          
          <ArticleList
            articles={isSearching ? searchResults : articles}
            onDeleteArticle={handleDeleteArticle}
            onAddTag={handleAddTag}
            onRemoveTag={handleRemoveTag}
            onViewArticle={handleViewArticle}
          />

          {/* 文章详情模态框 */}
          {selectedArticle && (
            <div className="article-detail-overlay" onClick={closeArticleDetail}>
              <div className="article-detail-modal" onClick={(e) => e.stopPropagation()}>
                <div className="article-detail-header">
                  <h1>{selectedArticle.title}</h1>
                  <button onClick={closeArticleDetail} className="close-btn">×</button>
                </div>
                <div className="article-detail-content">
                  {selectedArticle.image_url && (
                    <img src={selectedArticle.image_url} alt={selectedArticle.title} className="article-detail-image" />
                  )}
                  <div className="article-content">
                    {selectedArticle.content ? (
                      <div dangerouslySetInnerHTML={{ __html: selectedArticle.content }} />
                    ) : (
                      <p>{selectedArticle.excerpt}</p>
                    )}
                  </div>
                  {selectedArticle.tags && selectedArticle.tags.length > 0 && (
                    <div className="tags-section">
                      <strong>Tags: </strong>
                      <div className="modal-tags">
                        {selectedArticle.tags.map(tag => (
                          <span key={tag.id} className="tag-with-delete">
                            <span className="tag-name">{tag.name}</span>
                            <button 
                              className="tag-delete-btn"
                              onClick={async () => {
                                await handleRemoveTag(selectedArticle.id, tag.id);
                                // 更新模态框中的文章数据
                                const updatedArticle = await api.getArticle(selectedArticle.id);
                                setSelectedArticle(updatedArticle);
                              }}
                              title="删除标签"
                            >
                              ✕
                            </button>
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                  <div className="article-actions">
                    <a href={selectedArticle.url} target="_blank" rel="noopener noreferrer">
                      Read Original
                    </a>
                  </div>
                </div>
              </div>
            </div>
          )}
        </main>
      ) : (
        <main className="auth-required">
          <div className="auth-prompt">
            <h3>请先登录或注册以使用此服务</h3>
            <p>登录后您可以保存文章并按用户隔离管理您的内容。</p>
            <button 
              onClick={() => setShowAuthModal(true)} 
              className="auth-prompt-btn"
            >
              登录/注册
            </button>
          </div>
        </main>
      )}

      {/* 认证模态框 */}
      <AuthModal
        isOpen={showAuthModal}
        onClose={() => setShowAuthModal(false)}
        onLogin={handleLogin}
        onRegister={handleRegister}
        isLoading={authLoading}
      />
    </div>
  );
}

export default App;
