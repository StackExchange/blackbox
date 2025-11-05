'use client';

import { useEffect, useState } from 'react';

export function HeroSection() {
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    setVisible(true);
  }, []);

  return (
    <section className="relative min-h-screen flex items-center justify-center bg-gradient-to-br from-pink-50 via-rose-50 to-red-50 overflow-hidden">
      <div className="absolute inset-0 overflow-hidden">
        {[...Array(20)].map((_, i) => (
          <div
            key={i}
            className="absolute animate-float"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 5}s`,
              animationDuration: `${5 + Math.random() * 10}s`,
            }}
          >
            <span className="text-red-300 text-2xl opacity-30">❤️</span>
          </div>
        ))}
      </div>

      <div className={`relative z-10 text-center px-4 transition-all duration-2000 ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'}`}>
        <h1 className="font-playfair text-6xl md:text-8xl font-bold text-rose-600 mb-6 animate-fade-in">
          For My Love
        </h1>
        <p className="font-lato text-xl md:text-2xl text-gray-700 mb-8 animate-fade-in-delay">
          Every moment with you is a treasure
        </p>
        <div className="animate-bounce mt-12">
          <svg
            className="w-8 h-8 mx-auto text-rose-400"
            fill="none"
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path d="M19 14l-7 7m0 0l-7-7m7 7V3"></path>
          </svg>
        </div>
      </div>

      <style jsx>{`
        @keyframes float {
          0%, 100% {
            transform: translateY(0) rotate(0deg);
          }
          50% {
            transform: translateY(-100px) rotate(180deg);
          }
        }
        .animate-float {
          animation: float linear infinite;
        }
        @keyframes fade-in {
          from {
            opacity: 0;
            transform: translateY(20px);
          }
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
        .animate-fade-in {
          animation: fade-in 1.5s ease-out forwards;
        }
        .animate-fade-in-delay {
          animation: fade-in 1.5s ease-out 0.5s forwards;
          opacity: 0;
        }
      `}</style>
    </section>
  );
}
