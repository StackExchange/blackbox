'use client';

import { useState, useEffect } from 'react';
import { HeroSection } from '@/components/HeroSection';
import { StorySection } from '@/components/StorySection';
import { GallerySection } from '@/components/GallerySection';
import { ProposalSection } from '@/components/ProposalSection';

export default function Home() {
  const [showContent, setShowContent] = useState(false);

  useEffect(() => {
    setTimeout(() => setShowContent(true), 500);
  }, []);

  return (
    <main className={`min-h-screen transition-opacity duration-1000 ${showContent ? 'opacity-100' : 'opacity-0'}`}>
      <HeroSection />
      <StorySection />
      <GallerySection />
      <ProposalSection />
    </main>
  );
}
