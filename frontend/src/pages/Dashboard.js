import React, { useState, useEffect, useContext } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthContext from '../context/AuthContext';
import './Dashboard.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const Dashboard = () => {
  const [readings, setReadings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [marking, setMarking] = useState(false);
  const { token, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }

    fetchTodayReadings();
  }, [token, navigate]);

  const fetchTodayReadings = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_URL}/api/readings/today`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setReadings(response.data);
      setError('');
    } catch (err) {
      if (err.response?.status === 401) {
        logout();
        navigate('/login');
      } else {
        setError('Erro ao carregar leituras do dia');
      }
    } finally {
      setLoading(false);
    }
  };

  const markAsCompleted = async (period) => {
    try {
      setMarking(true);
      const response = await axios.post(
        `${API_URL}/api/readings/mark-completed`,
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
          <div className="loading">Carregando leituras do dia...</div>
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
          <p className="day-info">Dia {readings?.day_of_year} do ano</p>
        </div>

        {error && <div className="error-message">{error}</div>}

        <div className="readings-grid">
          {period === 'morning' || period === 'all' ? (
            <div className="reading-card">
              <h3>Antigo Testamento</h3>
              {readingPlan?.old_testament_ref ? (
                <>
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
                </>
              ) : (
                <p className="no-reading">Sem leitura para hoje</p>
              )}
            </div>
          ) : null}

          {period === 'morning' || period === 'all' ? (
            <div className="reading-card">
              <h3>Salmos</h3>
              {readingPlan?.psalms_ref ? (
                <>
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
                </>
              ) : (
                <p className="no-reading">Sem leitura para hoje</p>
              )}
            </div>
          ) : null}

          {period === 'evening' || period === 'all' ? (
            <div className="reading-card">
              <h3>Novo Testamento</h3>
              {readingPlan?.new_testament_ref ? (
                <>
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
                </>
              ) : (
                <p className="no-reading">Sem leitura para hoje</p>
              )}
            </div>
          ) : null}

          {period === 'evening' || period === 'all' ? (
            <div className="reading-card">
              <h3>Provérbios</h3>
              {readingPlan?.proverbs_ref ? (
                <>
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
                </>
              ) : (
                <p className="no-reading">Sem leitura para hoje</p>
              )}
            </div>
          ) : null}
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

