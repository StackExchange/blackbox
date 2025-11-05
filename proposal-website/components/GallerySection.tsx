'use client';

import { useEffect, useState, useRef } from 'react';

const memories = [
  {
    title: "Our Adventures",
    description: "Every journey with you is an adventure I cherish",
    color: "from-purple-400 to-pink-400"
  },
  {
    title: "Quiet Moments",
    description: "The simple times together mean the most",
    color: "from-blue-400 to-cyan-400"
  },
  {
    title: "Laughter & Joy",
    description: "Your smile lights up my entire world",
    color: "from-yellow-400 to-orange-400"
  },
  {
    title: "Growing Together",
    description: "Building a life and future side by side",
    color: "from-green-400 to-emerald-400"
  },
  {
    title: "Special Celebrations",
    description: "Every milestone is better with you",
    color: "from-red-400 to-rose-400"
  },
  {
    title: "Everyday Magic",
    description: "Finding beauty in the ordinary moments",
    color: "from-indigo-400 to-purple-400"
  }
];

export function GallerySection() {
  const [visible, setVisible] = useState(false);
  const sectionRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setVisible(true);
          }
        });
      },
      { threshold: 0.1 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <section ref={sectionRef} className="py-20 px-4 bg-gradient-to-br from-gray-50 to-rose-50">
      <div className="max-w-6xl mx-auto">
        <h2 className="font-playfair text-5xl md:text-6xl font-bold text-center text-gray-800 mb-4">
          Our Memories
        </h2>
        <p className="font-lato text-center text-gray-600 mb-16 text-lg">
          A collection of moments that made us who we are
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {memories.map((memory, index) => (
            <div
              key={index}
              className={`group relative overflow-hidden rounded-xl shadow-lg transition-all duration-700 hover:scale-105 hover:shadow-2xl ${
                visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'
              }`}
              style={{ transitionDelay: `${index * 100}ms` }}
            >
              <div className={`aspect-square bg-gradient-to-br ${memory.color} flex items-center justify-center p-8 relative overflow-hidden`}>
                <div className="absolute inset-0 bg-black opacity-0 group-hover:opacity-10 transition-opacity"></div>
                
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="text-white text-8xl opacity-20 group-hover:scale-110 transition-transform">
                    ❤️
                  </div>
                </div>

                <div className="relative z-10 text-center text-white">
                  <h3 className="font-playfair text-2xl font-bold mb-3 drop-shadow-lg">
                    {memory.title}
                  </h3>
                  <p className="font-lato text-sm opacity-90 drop-shadow">
                    {memory.description}
                  </p>
                </div>
              </div>

              <div className="absolute inset-0 border-4 border-white rounded-xl pointer-events-none"></div>
            </div>
          ))}
        </div>

        <div className="mt-12 text-center">
          <p className="font-lato text-gray-500 italic">
            Replace these colorful placeholders with your actual photos to make it even more special! 📸
          </p>
        </div>
      </div>
    </section>
  );
}
