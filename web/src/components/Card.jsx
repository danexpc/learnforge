export default function Card({ children, title, className = '' }) {
  return (
    <div className={`bg-white rounded-xl shadow-lg p-6 border border-gray-100 ${className}`}>
      {title && (
        <h3 className="text-xl font-bold text-gray-800 mb-4">{title}</h3>
      )}
      {children}
    </div>
  );
}

