'use client';

import { useEffect, useState, useRef } from 'react';

export function ProposalSection() {
  const [visible, setVisible] = useState(false);
  const [showQuestion, setShowQuestion] = useState(false);
  const [showButtons, setShowButtons] = useState(false);
  const [answered, setAnswered] = useState(false);
  const [hearts, setHearts] = useState<Array<{ id: number; x: number; y: number }>>([]);
  const sectionRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setVisible(true);
            setTimeout(() => setShowQuestion(true), 1000);
            setTimeout(() => setShowButtons(true), 2000);
          }
        });
      },
      { threshold: 0.3 }
    );

    if (sectionRef.current) {
      observer.observe(sectionRef.current);
    }

    return () => observer.disconnect();
  }, []);

  const handleYes = () => {
    setAnswered(true);
    createHeartExplosion();
  };

  const createHeartExplosion = () => {
    const newHearts = Array.from({ length: 50 }, (_, i) => ({
      id: Date.now() + i,
      x: Math.random() * 100,
      y: Math.random() * 100,
    }));
    setHearts(newHearts);

    setTimeout(() => {
      setHearts([]);
    }, 3000);
  };

  const handleNo = (e: React.MouseEvent<HTMLButtonElement>) => {
    const button = e.currentTarget;
    const maxX = window.innerWidth - button.offsetWidth - 20;
    const maxY = window.innerHeight - button.offsetHeight - 20;
    
    const randomX = Math.random() * maxX;
    const randomY = Math.random() * maxY;
    
    button.style.position = 'fixed';
    button.style.left = `${randomX}px`;
    button.style.top = `${randomY}px`;
  };

  return (
    <section 
      ref={sectionRef} 
      className="relative min-h-screen flex items-center justify-center bg-gradient-to-br from-rose-100 via-pink-100 to-red-100 overflow-hidden"
    >
      {hearts.map((heart) => (
        <div
          key={heart.id}
          className="absolute text-4xl animate-heart-float pointer-events-none"
          style={{
            left: `${heart.x}%`,
            top: `${heart.y}%`,
            animationDuration: `${2 + Math.random() * 2}s`,
          }}
        >
          ❤️
        </div>
      ))}

      <div className="relative z-10 text-center px-4 max-w-4xl">
        <div className={`transition-all duration-1000 ${visible ? 'opacity-100 scale-100' : 'opacity-0 scale-95'}`}>
          <div className="mb-12">
            <div className="text-8xl mb-8 animate-pulse">💍</div>
          </div>

          {!answered ? (
            <>
              <h2 className={`font-playfair text-4xl md:text-6xl font-bold text-rose-600 mb-8 transition-all duration-1000 ${showQuestion ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'}`}>
                You are my everything
              </h2>

              <p className={`font-lato text-xl md:text-2xl text-gray-700 mb-12 leading-relaxed transition-all duration-1000 delay-300 ${showQuestion ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'}`}>
                Every day with you is a gift. You make me laugh, you make me think, 
                you make me want to be a better person. I cannot imagine my life without you, 
                and I don&apos;t want to spend another day not being able to call you my forever.
              </p>

              <h1 className={`font-playfair text-5xl md:text-7xl font-bold text-red-600 mb-16 transition-all duration-1000 delay-500 ${showQuestion ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'}`}>
                Will You Marry Me?
              </h1>

              <div className={`flex gap-6 justify-center items-center transition-all duration-1000 delay-700 ${showButtons ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-10'}`}>
                <button
                  onClick={handleYes}
                  className="font-lato px-12 py-4 bg-gradient-to-r from-rose-500 to-red-500 text-white text-2xl font-bold rounded-full shadow-2xl hover:shadow-3xl hover:scale-110 transition-all duration-300 hover:from-rose-600 hover:to-red-600"
                >
                  Yes! 💕
                </button>

                <button
                  onClick={handleNo}
                  className="font-lato px-12 py-4 bg-gray-300 text-gray-600 text-2xl font-bold rounded-full shadow-lg hover:shadow-xl transition-all duration-300"
                >
                  No
                </button>
              </div>

              <p className="font-lato text-sm text-gray-500 mt-8 italic">
                (Hint: The &quot;No&quot; button is shy... it might run away! 😉)
              </p>
            </>
          ) : (
            <div className="animate-fade-in">
              <h2 className="font-playfair text-5xl md:text-7xl font-bold text-rose-600 mb-8">
                She Said Yes! 🎉
              </h2>
              <p className="font-lato text-2xl md:text-3xl text-gray-700 mb-8">
                I&apos;m the luckiest person in the world! ❤️
              </p>
              <div className="text-6xl animate-bounce">
                💑
              </div>
              <p className="font-lato text-xl text-gray-600 mt-12 italic">
                Forever starts now...
              </p>
            </div>
          )}
        </div>
      </div>

      <style jsx>{`
        @keyframes heart-float {
          0% {
            transform: translateY(0) scale(0) rotate(0deg);
            opacity: 1;
          }
          50% {
            opacity: 1;
          }
          100% {
            transform: translateY(-100vh) scale(1.5) rotate(360deg);
            opacity: 0;
          }
        }
        .animate-heart-float {
          animation: heart-float ease-out forwards;
        }
        @keyframes fade-in {
          from {
            opacity: 0;
            transform: scale(0.8);
          }
          to {
            opacity: 1;
            transform: scale(1);
          }
        }
        .animate-fade-in {
          animation: fade-in 1s ease-out forwards;
        }
      `}</style>
    </section>
  );
}
