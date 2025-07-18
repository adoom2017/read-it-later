import React, { useState, useEffect } from 'react';
import axios from 'axios';
import ArticleList from './components/ArticleList';
import './App.css';

function App() {
  const [url, setUrl] = useState('');
  const [articles, setArticles] = useState([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState(null);

  // Fetch articles on component mount
  useEffect(() => {
    const fetchArticles = async () => {
      try {
        const response = await axios.get('/api/articles');
        setArticles(response.data || []);
      } catch (err) {
        setError('Failed to fetch articles.');
        console.error(err);
      }
    };
    fetchArticles();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!url) return;

    setLoading(true);
    setError('');

    try {
      const response = await axios.post('/api/articles', { url });
      setArticles([response.data, ...articles]);
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
      await axios.delete(`/api/articles/${articleId}`);
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
      await axios.post(`/api/articles/${articleId}/tags`, { tag_name: tagName.trim() });
      // 重新获取文章列表以显示更新的标签
      const response = await axios.get('/api/articles');
      setArticles(response.data || []);
      setError('');
    } catch (err) {
      setError('Failed to add tag.');
      console.error(err);
    }
  };

  // 查看文章详情功能
  const handleViewArticle = async (articleId) => {
    try {
      const response = await axios.get(`/api/articles/${articleId}`);
      setSelectedArticle(response.data);
      setError('');
    } catch (err) {
      setError('Failed to fetch article details.');
      console.error(err);
    }
  };

  const closeArticleDetail = () => {
    setSelectedArticle(null);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>Read It Later</h1>
        <p>Save articles to read and summarize them.</p>
      </header>
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
        <ArticleList
          articles={articles}
          onDeleteArticle={handleDeleteArticle}
          onAddTag={handleAddTag}
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
                    {selectedArticle.tags.map(tag => (
                      <span key={tag.id} className="tag">{tag.name}</span>
                    ))}
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
    </div>
  );
}

export default App;
