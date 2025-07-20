import { useState, useEffect } from 'react';

const API_BASE_URL = 'http://localhost:8080';

// 获取存储的token
export const getToken = () => {
  return localStorage.getItem('auth_token');
};

// 存储token
export const setToken = (token) => {
  localStorage.setItem('auth_token', token);
};

// 删除token
export const removeToken = () => {
  localStorage.removeItem('auth_token');
};

// 检查是否已登录
export const isAuthenticated = () => {
  const token = getToken();
  if (!token) return false;
  
  try {
    // 简单检查token是否过期（这里可以更严格地验证JWT）
    const payload = JSON.parse(atob(token.split('.')[1]));
    return payload.exp * 1000 > Date.now();
  } catch {
    return false;
  }
};

// API请求工具函数
const apiRequest = async (url, options = {}) => {
  const token = getToken();
  
  const config = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
    ...options,
  };

  const response = await fetch(`${API_BASE_URL}${url}`, config);
  
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `HTTP ${response.status}`);
  }
  
  return response.json();
};

// 用户认证Hook
export const useAuth = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // 检查登录状态
  useEffect(() => {
    const checkAuth = async () => {
      if (isAuthenticated()) {
        try {
          const userData = await apiRequest('/api/user/profile');
          setUser(userData);
        } catch (err) {
          console.error('Failed to fetch user profile:', err);
          removeToken();
        }
      }
      setLoading(false);
    };

    checkAuth();
  }, []);

  // 登录
  const login = async (username, password) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiRequest('/api/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password }),
      });
      
      setToken(response.token);
      setUser(response.user);
      return { success: true };
    } catch (err) {
      setError(err.message);
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  };

  // 注册
  const register = async (username, email, password) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiRequest('/api/auth/register', {
        method: 'POST',
        body: JSON.stringify({ username, email, password }),
      });
      
      setToken(response.token);
      setUser(response.user);
      return { success: true };
    } catch (err) {
      setError(err.message);
      return { success: false, error: err.message };
    } finally {
      setLoading(false);
    }
  };

  // 登出
  const logout = () => {
    removeToken();
    setUser(null);
    setError(null);
  };

  return {
    user,
    loading,
    error,
    isAuthenticated: !!user,
    login,
    register,
    logout,
  };
};

// API调用工具函数
export const api = {
  // 文章相关API
  getArticles: () => apiRequest('/api/articles'),
  
  getArticle: (id) => apiRequest(`/api/articles/${id}`),
  
  createArticle: (articleData) => apiRequest('/api/articles', {
    method: 'POST',
    body: JSON.stringify(articleData),
  }),
  
  updateArticle: (id, articleData) => apiRequest(`/api/articles/${id}`, {
    method: 'PUT',
    body: JSON.stringify(articleData),
  }),
  
  deleteArticle: (id) => apiRequest(`/api/articles/${id}`, {
    method: 'DELETE',
  }),

  // 标签相关API
  getTags: () => apiRequest('/api/tags'),
  
  getArticlesByTag: (tagName) => apiRequest(`/api/tags/${tagName}/articles`),
  
  addTagToArticle: (articleId, tagName) => apiRequest(`/api/articles/${articleId}/tags`, {
    method: 'POST',
    body: JSON.stringify({ tag_name: tagName }),
  }),
  
  removeTagFromArticle: (articleId, tagId) => apiRequest(`/api/articles/${articleId}/tags/${tagId}`, {
    method: 'DELETE',
  }),

  // 搜索API
  searchArticles: (query, searchType) => {
    const params = new URLSearchParams();
    if (searchType === 'tag') {
      params.append('tag', query);
    } else {
      params.append('q', query);
    }
    return apiRequest(`/api/articles/search?${params}`);
  },
};

export default useAuth;
