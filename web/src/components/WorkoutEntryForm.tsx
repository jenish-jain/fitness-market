import React, { useState, useEffect } from 'react';
import { api } from '../services/api';

interface Exercise {
  id: number;
  ticker: string;
  name: string;
  category: string;
}

interface WorkoutEntryFormProps {
  onSuccess?: (entry: any) => void;
  onError?: (error: string) => void;
}

export const WorkoutEntryForm: React.FC<WorkoutEntryFormProps> = ({ onSuccess, onError }) => {
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    exercise_id: '',
    weight: '',
    reps: '',
    sets: '',
    notes: '',
    date: new Date().toISOString().split('T')[0]
  });
  const [result, setResult] = useState<{ score: number; is_pr: boolean } | null>(null);

  useEffect(() => {
    fetchExercises();
  }, []);

  const fetchExercises = async () => {
    setLoading(true);
    try {
      const response = await api.get('/exercises');
      setExercises(response.data.exercises || []);
    } catch (err) {
      onError?.('Failed to load exercises');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setResult(null);

    try {
      const payload = {
        exercise_id: parseInt(formData.exercise_id),
        weight: parseFloat(formData.weight),
        reps: parseInt(formData.reps),
        sets: parseInt(formData.sets),
        notes: formData.notes,
        date: formData.date
      };

      const response = await api.post('/entries', payload);
      const entry = response.data.entry;
      
      setResult({
        score: entry.score,
        is_pr: entry.is_pr
      });

      // Reset form
      setFormData({
        exercise_id: '',
        weight: '',
        reps: '',
        sets: '',
        notes: '',
        date: new Date().toISOString().split('T')[0]
      });

      onSuccess?.(entry);
    } catch (err: any) {
      const errorMsg = err.response?.data?.error || 'Failed to log workout entry';
      onError?.(errorMsg);
    } finally {
      setSubmitting(false);
    }
  };

  const isFormValid = formData.exercise_id && formData.weight && formData.reps && formData.sets;

  return (
    <div className="workout-entry-form">
      <h2 className="form-title">Log Workout Entry</h2>
      
      {result && (
        <div className={`result-banner ${result.is_pr ? 'pr-banner' : 'score-banner'}`}>
          {result.is_pr && (
            <div className="pr-celebration">
              🎉 New Personal Record! 🎉
            </div>
          )}
          <div className="score-display">
            Score: <span className="score-value">{result.score.toFixed(2)}</span>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="entry-form">
        <div className="form-group">
          <label htmlFor="exercise_id">Exercise</label>
          <select
            id="exercise_id"
            name="exercise_id"
            value={formData.exercise_id}
            onChange={handleChange}
            required
            disabled={loading}
            className="form-control"
          >
            <option value="">Select an exercise</option>
            {exercises.map(ex => (
              <option key={ex.id} value={ex.id}>
                {ex.name} ({ex.ticker})
              </option>
            ))}
          </select>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label htmlFor="weight">Weight (kg)</label>
            <input
              type="number"
              id="weight"
              name="weight"
              value={formData.weight}
              onChange={handleChange}
              required
              min="0"
              step="0.5"
              className="form-control"
              placeholder="0.0"
            />
          </div>

          <div className="form-group">
            <label htmlFor="reps">Reps</label>
            <input
              type="number"
              id="reps"
              name="reps"
              value={formData.reps}
              onChange={handleChange}
              required
              min="1"
              className="form-control"
              placeholder="0"
            />
          </div>

          <div className="form-group">
            <label htmlFor="sets">Sets</label>
            <input
              type="number"
              id="sets"
              name="sets"
              value={formData.sets}
              onChange={handleChange}
              required
              min="1"
              className="form-control"
              placeholder="0"
            />
          </div>
        </div>

        <div className="form-group">
          <label htmlFor="date">Date</label>
          <input
            type="date"
            id="date"
            name="date"
            value={formData.date}
            onChange={handleChange}
            className="form-control"
          />
        </div>

        <div className="form-group">
          <label htmlFor="notes">Notes (optional)</label>
          <textarea
            id="notes"
            name="notes"
            value={formData.notes}
            onChange={handleChange}
            className="form-control"
            rows={3}
            placeholder="Add any notes about this workout..."
          />
        </div>

        <button
          type="submit"
          disabled={!isFormValid || submitting}
          className="submit-btn"
        >
          {submitting ? 'Logging...' : 'Log Workout'}
        </button>
      </form>

      <style>{`
        .workout-entry-form {
          max-width: 600px;
          margin: 0 auto;
          padding: 20px;
          background: #fff;
          border-radius: 12px;
          box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
        }

        .form-title {
          font-size: 1.5rem;
          font-weight: 600;
          margin-bottom: 20px;
          color: #333;
          text-align: center;
        }

        .result-banner {
          padding: 16px;
          border-radius: 8px;
          margin-bottom: 20px;
          text-align: center;
        }

        .pr-banner {
          background: linear-gradient(135deg, #ffd700 0%, #ffb347 100%);
          color: #333;
        }

        .score-banner {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: #fff;
        }

        .pr-celebration {
          font-size: 1.25rem;
          font-weight: 700;
          margin-bottom: 8px;
        }

        .score-display {
          font-size: 1.1rem;
        }

        .score-value {
          font-weight: 700;
          font-size: 1.3rem;
        }

        .entry-form {
          display: flex;
          flex-direction: column;
          gap: 16px;
        }

        .form-group {
          display: flex;
          flex-direction: column;
          gap: 6px;
          flex: 1;
        }

        .form-group label {
          font-size: 0.875rem;
          font-weight: 500;
          color: #555;
        }

        .form-control {
          padding: 12px;
          border: 1px solid #ddd;
          border-radius: 8px;
          font-size: 1rem;
          transition: border-color 0.2s, box-shadow 0.2s;
        }

        .form-control:focus {
          outline: none;
          border-color: #667eea;
          box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .form-row {
          display: flex;
          gap: 12px;
        }

        .submit-btn {
          padding: 14px 24px;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
          color: #fff;
          border: none;
          border-radius: 8px;
          font-size: 1rem;
          font-weight: 600;
          cursor: pointer;
          transition: transform 0.2s, box-shadow 0.2s;
        }

        .submit-btn:hover:not(:disabled) {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
        }

        .submit-btn:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }

        /* Mobile responsiveness */
        @media (max-width: 768px) {
          .workout-entry-form {
            padding: 16px;
            margin: 10px;
            border-radius: 8px;
          }

          .form-title {
            font-size: 1.25rem;
          }

          .form-row {
            flex-direction: column;
            gap: 16px;
          }

          .form-control {
            padding: 14px;
            font-size: 16px; /* Prevents zoom on iOS */
          }

          .submit-btn {
            padding: 16px;
            font-size: 1.1rem;
          }

          .pr-celebration {
            font-size: 1.1rem;
          }

          .score-display {
            font-size: 1rem;
          }

          .score-value {
            font-size: 1.2rem;
          }
        }

        @media (max-width: 480px) {
          .workout-entry-form {
            padding: 12px;
            margin: 8px;
          }

          .form-group label {
            font-size: 0.8rem;
          }
        }
      `}</style>
    </div>
  );
};

export default WorkoutEntryForm;
