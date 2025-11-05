'use client';

import { useEffect, useState, useRef } from 'react';

const milestones = [
  {
    title: "The First Hello",
    description: "The day our eyes met and my world changed forever. I knew from that moment you were special.",
    date: "Day One"
  },
  {
    title: "Our First Date",
    description: "Nervous butterflies, endless conversation, and the beginning of something beautiful.",
    date: "The Beginning"
  },
  {
    title: "Falling Deeper",
    description: "Every laugh, every smile, every moment together made me fall more in love with you.",
    date: "Every Day Since"
  },
  {
    title: "Building Dreams",
    description: "Planning our future, sharing our hopes, and realizing I want to spend forever with you.",
    date: "Our Journey"
  }
];

export function StorySection() {
  const [visibleItems, setVisibleItems] = useState<number[]>([]);
  const sectionRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            milestones.forEach((_, index) => {
              setTimeout(() => {
                setVisibleItems((prev) => [...prev, index]);
              }, index * 300);
            });
          }
        });
      },
      { threshold: 0.2 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  return (
    <section ref={sectionRef} className="py-20 px-4 bg-white">
      <div className="max-w-4xl mx-auto">
        <h2 className="font-playfair text-5xl md:text-6xl font-bold text-center text-gray-800 mb-16">
          Our Story
        </h2>

        <div className="relative">
          <div className="absolute left-1/2 transform -translate-x-1/2 h-full w-1 bg-gradient-to-b from-rose-300 to-red-400"></div>

          {milestones.map((milestone, index) => (
            <div
              key={index}
              className={`relative mb-16 transition-all duration-1000 ${
                visibleItems.includes(index)
                  ? 'opacity-100 translate-x-0'
                  : index % 2 === 0
                  ? 'opacity-0 -translate-x-10'
                  : 'opacity-0 translate-x-10'
              }`}
            >
              <div className={`flex items-center ${index % 2 === 0 ? 'flex-row' : 'flex-row-reverse'}`}>
                <div className={`w-1/2 ${index % 2 === 0 ? 'pr-8 text-right' : 'pl-8 text-left'}`}>
                  <div className="bg-gradient-to-br from-rose-50 to-pink-50 p-6 rounded-lg shadow-lg hover:shadow-xl transition-shadow">
                    <h3 className="font-playfair text-2xl font-bold text-rose-600 mb-2">
                      {milestone.title}
                    </h3>
                    <p className="font-lato text-gray-600 mb-3">
                      {milestone.description}
                    </p>
                    <span className="font-lato text-sm text-rose-400 font-semibold">
                      {milestone.date}
                    </span>
                  </div>
                </div>

                <div className="absolute left-1/2 transform -translate-x-1/2 w-6 h-6 bg-rose-500 rounded-full border-4 border-white shadow-lg z-10">
                  <div className="absolute inset-0 bg-rose-500 rounded-full animate-ping opacity-75"></div>
                </div>

                <div className="w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
