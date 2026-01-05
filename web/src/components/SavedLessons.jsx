import { useState, useEffect } from 'preact/hooks';
import { getSavedLessons, deleteLesson } from '../utils/storage';
import { getResult } from '../utils/api';
import Card from './Card';
import Button from './Button';

export default function SavedLessons({ onLoadLesson }) {
  const [savedLessons, setSavedLessons] = useState([]);
  const [loadingIds, setLoadingIds] = useState(new Set());
  const [errorIds, setErrorIds] = useState(new Set());

  useEffect(() => {
    loadSavedLessons();
  }, []);

  const loadSavedLessons = () => {
    const lessons = getSavedLessons();
    setSavedLessons(lessons);
    
    lessons.forEach(lesson => {
      verifyLessonExists(lesson.id);
    });
  };

  const verifyLessonExists = async (id) => {
    setLoadingIds(prev => new Set(prev).add(id));
    setErrorIds(prev => {
      const newSet = new Set(prev);
      newSet.delete(id);
      return newSet;
    });

    try {
      await getResult(id);
    } catch (err) {
      if (err.message.includes('not found') || err.message.includes('404')) {
        setErrorIds(prev => new Set(prev).add(id));
      }
    } finally {
      setLoadingIds(prev => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
    }
  };

  const handleLoad = async (lesson) => {
    if (onLoadLesson) {
      onLoadLesson(lesson.id);
    }
  };

  const handleDelete = (id, e) => {
    e.stopPropagation();
    if (confirm('Are you sure you want to delete this lesson?')) {
      if (deleteLesson(id)) {
        setSavedLessons(prev => prev.filter(l => l.id !== id));
        setErrorIds(prev => {
          const newSet = new Set(prev);
          newSet.delete(id);
          return newSet;
        });
      }
    }
  };

  if (savedLessons.length === 0) {
    return null;
  }

  return (
    <Card title="📚 Saved Lessons">
      <div className="space-y-3">
        {savedLessons.map((lesson) => {
          const isLoading = loadingIds.has(lesson.id);
          const hasError = errorIds.has(lesson.id);

          return (
            <div
              key={lesson.id}
              className={`p-4 border rounded-lg transition-all ${
                hasError
                  ? 'border-red-200 bg-red-50'
                  : 'border-gray-200 bg-white hover:border-blue-300 hover:shadow-md cursor-pointer'
              }`}
              onClick={() => !hasError && handleLoad(lesson)}
            >
              <div className="flex items-start justify-between">
                <div className="flex-1 min-w-0">
                  <h3 className="font-semibold text-gray-800 truncate">
                    {lesson.topic || 'Untitled Lesson'}
                  </h3>
                  <div className="mt-1 flex items-center gap-3 text-sm text-gray-600">
                    <span className="capitalize">{lesson.mode || 'lesson'}</span>
                    <span>•</span>
                    <span className="capitalize">{lesson.level || 'beginner'}</span>
                    {lesson.createdAt && (
                      <>
                        <span>•</span>
                        <span>{new Date(lesson.createdAt).toLocaleDateString()}</span>
                      </>
                    )}
                  </div>
                  {hasError && (
                    <p className="mt-2 text-sm text-red-600">
                      ⚠️ This lesson no longer exists on the server
                    </p>
                  )}
                </div>
                <div className="flex items-center gap-2 ml-4">
                  {isLoading && (
                    <div className="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin"></div>
                  )}
                  <button
                    onClick={(e) => handleDelete(lesson.id, e)}
                    className="p-1 text-gray-400 hover:text-red-600 transition-colors"
                    title="Delete lesson"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </Card>
  );
}

