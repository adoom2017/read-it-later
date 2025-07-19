import React, { useState } from 'react';

const ArticleList = ({ articles, onDeleteArticle, onAddTag, onRemoveTag, onViewArticle }) => {
  const [tagInputs, setTagInputs] = useState({});

  const handleTagInputChange = (articleId, value) => {
    setTagInputs(prev => ({
      ...prev,
      [articleId]: value
    }));
  };

  const handleAddTag = async (articleId) => {
    const tagName = tagInputs[articleId];
    if (!tagName || !tagName.trim()) return;

    await onAddTag(articleId, tagName.trim());
    setTagInputs(prev => ({
      ...prev,
      [articleId]: ''
    }));
  };

  const handleKeyPress = (e, articleId) => {
    if (e.key === 'Enter') {
      handleAddTag(articleId);
    }
  };

  if (!articles.length) {
    return <p>No articles saved yet. Paste a link above to get started!</p>;
  }

  return (
    <div className="article-list">
      {articles.map(article => (
        <div key={article.id} className="article-card">
          {article.image_url && <img src={article.image_url} alt={article.title} />}
          <div className="article-card-content">
            <div className="title-with-tags">
              <h2>{article.title}</h2>
              {/* 在标题后面显示标签 */}
              {article.tags && article.tags.length > 0 && (
                <div className="inline-tags">
                  {article.tags.map(tag => (
                    <span key={tag.id} className="tag-with-delete">
                      <span className="tag-name">{tag.name}</span>
                      <button 
                        className="tag-delete-btn"
                        onClick={() => onRemoveTag(article.id, tag.id)}
                        title="删除标签"
                      >
                        ✕
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </div>
            <p>{article.excerpt}</p>

            {/* 添加标签区域 */}
            <div className="add-tag-section">
              <input
                type="text"
                placeholder="Add a tag..."
                value={tagInputs[article.id] || ''}
                onChange={(e) => handleTagInputChange(article.id, e.target.value)}
                onKeyPress={(e) => handleKeyPress(e, article.id)}
                className="tag-input"
              />
              <button
                onClick={() => handleAddTag(article.id)}
                className="add-tag-btn"
              >
                Add Tag
              </button>
            </div>

            {/* 文章操作按钮 */}
            <div className="article-actions">
              <a href={article.url} target="_blank" rel="noopener noreferrer">Read Original</a>
              <button
                onClick={() => onViewArticle(article.id)}
                className="view-btn"
              >
                View Details
              </button>
              <button
                onClick={() => onDeleteArticle(article.id)}
                className="delete-btn"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ArticleList;
