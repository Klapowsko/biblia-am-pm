import React, { useState, useEffect, useContext, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthContext from '../context/AuthContext';
import './Progress.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api';

const Progress = () => {
  const [progress, setProgress] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const { token, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const fetchProgress = useCallback(async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_URL}/progress`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setProgress(response.data);
      setError('');
    } catch (err) {
      if (err.response?.status === 401) {
        logout();
        navigate('/login');
      } else {
        setError('Erro ao carregar progresso');
      }
    } finally {
      setLoading(false);
    }
  }, [token, logout, navigate]);

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }

    fetchProgress();
  }, [token, navigate, fetchProgress]);

  const getCompletedDays = () => {
    return progress.filter(
      (p) => p.morning_completed && p.evening_completed
    ).length;
  };

  const getTotalDays = () => {
    return progress.length;
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="progress-container">
        <div className="container">
          <div className="loading">Carregando progresso...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="progress-container">
      <header className="dashboard-header">
        <div className="container">
          <h1>Bíblia AM/PM</h1>
          <div className="header-actions">
            <button
              className="btn btn-secondary"
              onClick={() => navigate('/')}
            >
              Dashboard
            </button>
            <button className="btn btn-secondary" onClick={logout}>
              Sair
            </button>
          </div>
        </div>
      </header>

      <div className="container">
        <div className="progress-header">
          <h2>Meu Progresso</h2>
          <div className="stats">
            <div className="stat-card">
              <div className="stat-value">{getCompletedDays()}</div>
              <div className="stat-label">Dias Completos</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">{getTotalDays()}</div>
              <div className="stat-label">Total de Dias</div>
            </div>
            <div className="stat-card">
              <div className="stat-value">
                {getTotalDays() > 0
                  ? Math.round((getCompletedDays() / getTotalDays()) * 100)
                  : 0}
                %
              </div>
              <div className="stat-label">Conclusão</div>
            </div>
          </div>
        </div>

        {error && <div className="error-message">{error}</div>}

        {progress.length === 0 ? (
          <div className="empty-state">
            <p>Você ainda não tem progresso registrado.</p>
            <p>Comece lendo as leituras do dia no dashboard!</p>
          </div>
        ) : (
          <div className="progress-list">
            <h3>Histórico de Leitura</h3>
            <div className="progress-table">
              <div className="table-header">
                <div>Data</div>
                <div>Manhã</div>
                <div>Noite</div>
                <div>Status</div>
              </div>
              {progress.map((item) => (
                <div
                  key={item.id}
                  className={`table-row ${
                    item.morning_completed && item.evening_completed
                      ? 'completed'
                      : ''
                  }`}
                >
                  <div>{formatDate(item.date)}</div>
                  <div>
                    {item.morning_completed ? (
                      <span className="check">✓</span>
                    ) : (
                      <span className="cross">✗</span>
                    )}
                  </div>
                  <div>
                    {item.evening_completed ? (
                      <span className="check">✓</span>
                    ) : (
                      <span className="cross">✗</span>
                    )}
                  </div>
                  <div>
                    {item.morning_completed && item.evening_completed ? (
                      <span className="status-badge completed">Completo</span>
                    ) : (
                      <span className="status-badge pending">Pendente</span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Progress;

