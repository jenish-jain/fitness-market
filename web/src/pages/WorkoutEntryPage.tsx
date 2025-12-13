import React, { useState } from 'react';
import { WorkoutEntryForm } from '../components/WorkoutEntryForm';

export const WorkoutEntryPage: React.FC = () => {
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

  const handleSuccess = (entry: any) => {
    const prText = entry.is_pr ? ' 🎉 New PR!' : '';
    setMessage({
      type: 'success',
      text: `Workout logged successfully! Score: ${entry.score.toFixed(2)}${prText}`
    });
    setTimeout(() => setMessage(null), 5000);
  };

  const handleError = (error: string) => {
    setMessage({ type: 'error', text: error });
    setTimeout(() => setMessage(null), 5000);
  };

  return (
    <div className="workout-entry-page">
      <div className="page-container">
        {message && (
          <div className={`message ${message.type}`}>
            {message.text}
          </div>
        )}
        
        <WorkoutEntryForm
          onSuccess={handleSuccess}
          onError={handleError}
        />
      </div>

      <style>{`
        .workout-entry-page {
          min-height: 100vh;
          background: #f5f7fa;
          padding: 20px;
        }

        .page-container {
          max-width: 800px;
          margin: 0 auto;
        }

        .message {
          padding: 12px 16px;
          border-radius: 8px;
          margin-bottom: 20px;
          font-weight: 500;
          text-align: center;
        }

        .message.success {
          background: #d4edda;
          color: #155724;
          border: 1px solid #c3e6cb;
        }

        .message.error {
          background: #f8d7da;
          color: #721c24;
          border: 1px solid #f5c6cb;
        }

        /* Mobile responsiveness */
        @media (max-width: 768px) {
          .workout-entry-page {
            padding: 10px;
          }

          .message {
            font-size: 0.9rem;
            padding: 10px 12px;
          }
        }
      `}</style>
    </div>
  );
};

export default WorkoutEntryPage;
