import { useState, useEffect } from 'react';

interface TypewriterTextProps {
  text: string;
  className?: string;
  typeSpeed?: number;
  deleteSpeed?: number;
  pauseAfterType?: number;
  pauseAfterDelete?: number;
}

type Phase = 'typing' | 'pause-typed' | 'deleting' | 'pause-deleted';

export function TypewriterText({
  text,
  className = '',
  typeSpeed = 120,
  deleteSpeed = 60,
  pauseAfterType = 2500,
  pauseAfterDelete = 1000,
}: TypewriterTextProps) {
  const [displayText, setDisplayText] = useState('');
  const [phase, setPhase] = useState<Phase>('typing');
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    let timeoutId: number;

    const runCycle = () => {
      switch (phase) {
        case 'typing':
          if (currentIndex < text.length) {
            timeoutId = setTimeout(() => {
              setDisplayText(text.slice(0, currentIndex + 1));
              setCurrentIndex(currentIndex + 1);
            }, typeSpeed);
          } else {
            // 打字完成，进入暂停状态
            timeoutId = setTimeout(() => {
              setPhase('pause-typed');
            }, pauseAfterType);
          }
          break;

        case 'pause-typed':
          // 暂停后开始删除
          timeoutId = setTimeout(() => {
            setPhase('deleting');
          }, 0);
          break;

        case 'deleting':
          if (currentIndex > 0) {
            timeoutId = setTimeout(() => {
              setDisplayText(text.slice(0, currentIndex - 1));
              setCurrentIndex(currentIndex - 1);
            }, deleteSpeed);
          } else {
            // 删除完成，进入暂停状态
            timeoutId = setTimeout(() => {
              setPhase('pause-deleted');
            }, pauseAfterDelete);
          }
          break;

        case 'pause-deleted':
          // 暂停后重新开始打字
          timeoutId = setTimeout(() => {
            setPhase('typing');
          }, 0);
          break;
      }
    };

    runCycle();

    // Cleanup：防止内存泄漏
    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [phase, currentIndex, text, typeSpeed, deleteSpeed, pauseAfterType, pauseAfterDelete]);

  // 文本改变时重置状态
  useEffect(() => {
    setDisplayText('');
    setCurrentIndex(0);
    setPhase('typing');
  }, [text]);

  return (
    <span className={className}>
      {displayText}
      <span className="animate-cursor">|</span>
      <style>{`
        @keyframes cursor-blink {
          0%, 49% { opacity: 1; }
          50%, 100% { opacity: 0; }
        }
        .animate-cursor {
          animation: cursor-blink 1s infinite;
        }
      `}</style>
    </span>
  );
}
