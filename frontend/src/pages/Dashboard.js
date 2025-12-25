import React, { useState, useEffect, useContext, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthContext from '../context/AuthContext';
import './Dashboard.css';

const API_URL = process.env.REACT_APP_API_URL || '/api';

const Dashboard = () => {
  const [readings, setReadings] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [marking, setMarking] = useState(false);
  const [catechism, setCatechism] = useState(null);
  const [catechismLoading, setCatechismLoading] = useState(true);
  const [catechismError, setCatechismError] = useState('');
  const [markingCatechism, setMarkingCatechism] = useState(false);
  const [showAnswer, setShowAnswer] = useState(false);
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

  const fetchCatechism = useCallback(async () => {
    try {
      setCatechismLoading(true);
      const response = await axios.get(`${API_URL}/catechism/current`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      setCatechism(response.data);
      setCatechismError('');
    } catch (err) {
      console.error('Error fetching catechism:', err);
      if (err.response?.status === 401) {
        logout();
        navigate('/login');
      } else if (err.response?.status === 404) {
        setCatechismError('Catecismo não populado. Por favor, popule o catecismo primeiro.');
      } else {
        setCatechismError('Erro ao carregar catecismo');
      }
    } finally {
      setCatechismLoading(false);
    }
  }, [token, logout, navigate]);

  useEffect(() => {
    if (!token) {
      navigate('/login');
      return;
    }

    fetchTodayReadings();
    fetchCatechism();
  }, [token, navigate, fetchTodayReadings, fetchCatechism]);

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

  const markCatechismAsCompleted = async () => {
    try {
      setMarkingCatechism(true);
      await axios.post(
        `${API_URL}/catechism/mark-completed`,
        {},
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      // Refresh catechism data
      await fetchCatechism();
    } catch (err) {
      setCatechismError('Erro ao marcar catecismo como concluído');
    } finally {
      setMarkingCatechism(false);
    }
  };

  const getWeekProgressCount = () => {
    if (!catechism?.week_progress) return 0;
    return catechism.week_progress.filter(p => p.completed).length;
  };

  const isTodayCompleted = () => {
    if (!catechism?.week_progress) return false;
    const today = new Date().toISOString().split('T')[0];
    const todayProgress = catechism.week_progress.find(p => {
      // Handle both date string formats (YYYY-MM-DD or ISO string)
      const progressDate = typeof p.date === 'string' 
        ? p.date.split('T')[0] 
        : new Date(p.date).toISOString().split('T')[0];
      return progressDate === today;
    });
    return todayProgress?.completed || false;
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
          <p className="plan-name">Plano: {readings?.plan_name || "Robert Murray M'Cheyne"}</p>
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

        {/* Catecismo Section */}
        <div className="catechism-section">
          <div className="catechism-header">
            <h2>📜 Catecismo Maior de Westminster</h2>
            {catechism && (
              <p className="catechism-info">
                Pergunta {catechism.question_number} de {catechism.total_questions}
              </p>
            )}
          </div>

          {catechismLoading && (
            <div className="catechism-loading">
              Carregando catecismo...
            </div>
          )}

          {catechismError && !catechismLoading && (
            <div className="catechism-error">
              {catechismError}
            </div>
          )}

          {catechism && !catechismLoading && (
            <div className="catechism-card">
              <div className="catechism-question">
                <h3>Pergunta {catechism.question?.question_number}</h3>
                <p className="question-text">{catechism.question?.question_text}</p>
              </div>

              <div className="catechism-answer">
                <button
                  className="btn-toggle-answer"
                  onClick={() => setShowAnswer(!showAnswer)}
                >
                  {showAnswer ? 'Ocultar Resposta' : 'Mostrar Resposta'}
                </button>
                {showAnswer && (
                  <div className="answer-content">
                    <p>{catechism.question?.answer_text}</p>
                  </div>
                )}
              </div>

              <div className="catechism-progress-info">
                <div className="week-progress">
                  <span>Progresso da Semana:</span>
                  <span className="progress-count">
                    {getWeekProgressCount()} / 7 dias
                  </span>
                </div>
                <div className="next-question-info">
                  <span>Próxima pergunta em: {new Date(catechism.next_question_date).toLocaleDateString('pt-BR')}</span>
                </div>
              </div>

              <button
                className={`btn ${isTodayCompleted() ? 'btn-success' : 'btn-primary'}`}
                onClick={markCatechismAsCompleted}
                disabled={markingCatechism || isTodayCompleted()}
              >
                {isTodayCompleted()
                  ? '✓ Lido hoje'
                  : 'Marcar como lido hoje'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

