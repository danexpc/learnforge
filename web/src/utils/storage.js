const STORAGE_KEY = 'learnforge_saved_lessons';

export function getSavedLessons() {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return [];
    return JSON.parse(stored);
  } catch (err) {
    console.error('Failed to load saved lessons:', err);
    return [];
  }
}

export function saveLesson(lessonData) {
  try {
    const lessons = getSavedLessons();
    const existingIndex = lessons.findIndex(l => l.id === lessonData.id);
    
    if (existingIndex >= 0) {
      lessons[existingIndex] = lessonData;
    } else {
      lessons.unshift(lessonData);
    }
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(lessons));
    return true;
  } catch (err) {
    console.error('Failed to save lesson:', err);
    return false;
  }
}

export function deleteLesson(id) {
  try {
    const lessons = getSavedLessons();
    const filtered = lessons.filter(l => l.id !== id);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(filtered));
    return true;
  } catch (err) {
    console.error('Failed to delete lesson:', err);
    return false;
  }
}

export function clearAllLessons() {
  try {
    localStorage.removeItem(STORAGE_KEY);
    return true;
  } catch (err) {
    console.error('Failed to clear lessons:', err);
    return false;
  }
}

