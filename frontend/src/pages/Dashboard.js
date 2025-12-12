import React, { useState, useEffect, useContext, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthContext from '../context/AuthContext';
import './Dashboard.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081/api';

const Dashboard = () => {
  const [readings, setReadings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [marking, setMarking] = useState(false);
  const { token, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  const fetchTodayReadings = useCallback(async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_URL}/readings/today`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setReadings(response.data);
      setError('');
    } catch (err) {
      console.error('Error fetching readings:', err);
      if (err.response?.status === 401) {
        logout();
        navigate('/login');
      } else {
        setError('Erro ao carregar leituras do dia');
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

    fetchTodayReadings();
  }, [token, navigate, fetchTodayReadings]);

  const markAsCompleted = async (period) => {
    try {
      setMarking(true);
      const response = await axios.post(
        `${API_URL}/readings/mark-completed`,
        { period },
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      // Update local state
      setReadings((prev) => ({
        ...prev,
        progress: response.data,
      }));
    } catch (err) {
      setError('Erro ao marcar leitura como concluída');
    } finally {
      setMarking(false);
    }
  };

  const getPeriodLabel = (period) => {
    switch (period) {
      case 'morning':
        return 'Manhã';
      case 'evening':
        return 'Noite';
      default:
        return 'Todo o dia';
    }
  };

  const getPeriodIcon = (period) => {
    switch (period) {
      case 'morning':
        return '☀️';
      case 'evening':
        return '🌙';
      default:
        return '📖';
    }
  };

  if (loading) {
    return (
      <div className="dashboard-container">
        <div className="container">
          <div className="loading">
            <div style={{ fontSize: '48px', marginBottom: '16px' }}>📖</div>
            Carregando leituras do dia...
          </div>
        </div>
      </div>
    );
  }

  if (error && !readings) {
    return (
      <div className="dashboard-container">
        <div className="container">
          <div className="error">{error}</div>
        </div>
      </div>
    );
  }

  const { period, readings: readingPlan, progress } = readings || {};

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <div className="container">
          <h1>Bíblia AM/PM</h1>
          <div className="header-actions">
            <button
              className="btn btn-secondary"
              onClick={() => navigate('/progress')}
            >
              Progresso
            </button>
            <button className="btn btn-secondary" onClick={logout}>
              Sair
            </button>
          </div>
        </div>
      </header>

      <div className="container">
        <div className="period-banner">
          <span className="period-icon">{getPeriodIcon(period)}</span>
          <h2>Leituras de {getPeriodLabel(period)}</h2>
          <p className="day-info">Dia {readings?.day_of_year} de 365</p>
        </div>

        {error && <div className="error-message">{error}</div>}

        <div className="readings-grid">
          {/* MANHÃ: Leituras (AT + NT) */}
          {(period === 'morning' || period === 'all') && readingPlan?.old_testament_ref && (
            <div className="reading-card">
              <h3>Leituras da Manhã</h3>
              <p className="reading-ref">{readingPlan.old_testament_ref}</p>
              <button
                className={`btn ${
                  progress?.morning_completed ? 'btn-success' : 'btn-primary'
                }`}
                onClick={() => markAsCompleted('morning')}
                disabled={marking || progress?.morning_completed}
              >
                {progress?.morning_completed
                  ? '✓ Concluído'
                  : 'Marcar como lido'}
              </button>
            </div>
          )}

          {/* Salmos - apenas se houver leitura */}
          {readingPlan?.psalms_ref && (
            <div className="reading-card">
              <h3>Salmos</h3>
              <p className="reading-ref">{readingPlan.psalms_ref}</p>
              <button
                className={`btn ${
                  progress?.morning_completed ? 'btn-success' : 'btn-primary'
                }`}
                onClick={() => markAsCompleted('morning')}
                disabled={marking || progress?.morning_completed}
              >
                {progress?.morning_completed
                  ? '✓ Concluído'
                  : 'Marcar como lido'}
              </button>
            </div>
          )}

          {/* NOITE: Leituras (AT + NT) */}
          {(period === 'evening' || period === 'all') && readingPlan?.new_testament_ref && (
            <div className="reading-card">
              <h3>Leituras da Noite</h3>
              <p className="reading-ref">{readingPlan.new_testament_ref}</p>
              <button
                className={`btn ${
                  progress?.evening_completed ? 'btn-success' : 'btn-primary'
                }`}
                onClick={() => markAsCompleted('evening')}
                disabled={marking || progress?.evening_completed}
              >
                {progress?.evening_completed
                  ? '✓ Concluído'
                  : 'Marcar como lido'}
              </button>
            </div>
          )}

          {/* Provérbios - apenas se houver leitura */}
          {readingPlan?.proverbs_ref && (
            <div className="reading-card">
              <h3>Provérbios</h3>
              <p className="reading-ref">{readingPlan.proverbs_ref}</p>
              <button
                className={`btn ${
                  progress?.evening_completed ? 'btn-success' : 'btn-primary'
                }`}
                onClick={() => markAsCompleted('evening')}
                disabled={marking || progress?.evening_completed}
              >
                {progress?.evening_completed
                  ? '✓ Concluído'
                  : 'Marcar como lido'}
              </button>
            </div>
          )}
        </div>

        <div className="progress-summary">
          <h3>Resumo do Dia</h3>
          <div className="progress-items">
            <div className="progress-item">
              <span>Manhã:</span>
              <span className={progress?.morning_completed ? 'completed' : 'pending'}>
                {progress?.morning_completed ? '✓ Concluído' : 'Pendente'}
              </span>
            </div>
            <div className="progress-item">
              <span>Noite:</span>
              <span className={progress?.evening_completed ? 'completed' : 'pending'}>
                {progress?.evening_completed ? '✓ Concluído' : 'Pendente'}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

