import React, { useState } from 'react';

const SearchBar = ({ onSearch, onClear }) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [searchType, setSearchType] = useState('title'); // 'title' or 'tag'

  const handleSubmit = (e) => {
    e.preventDefault();
    if (searchQuery.trim()) {
      onSearch(searchQuery.trim(), searchType);
    }
  };

  const handleClear = () => {
    setSearchQuery('');
    onClear();
  };

  return (
    <div className="search-container">
      <form onSubmit={handleSubmit} className="search-form">
        <div className="search-input-group">
          <select 
            value={searchType} 
            onChange={(e) => setSearchType(e.target.value)}
            className="search-type-select"
          >
            <option value="title">按标题搜索</option>
            <option value="tag">按标签搜索</option>
          </select>
          
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder={searchType === 'title' ? '搜索文章标题...' : '搜索标签...'}
            className="search-input"
          />
          
          <button type="submit" className="search-btn">
            🔍 搜索
          </button>
          
          {searchQuery && (
            <button type="button" onClick={handleClear} className="clear-btn">
              ✕ 清除
            </button>
          )}
        </div>
      </form>
    </div>
  );
};

export default SearchBar;
