declare module 'asciinema-player' {
  interface PlayerCreateOptions {
    cols?: number;
    rows?: number;
    autoPlay?: boolean;
    preload?: boolean;
    loop?: boolean;
    startAt?: number;
    speed?: number;
    idleTimeLimit?: number;
    theme?: string;
    poster?: string;
    fit?: 'width' | 'height' | 'both' | 'none' | false;
    controls?: boolean;
    title?: true | false | 'auto';
    markers?: any[];
    pauseOnMarkers?: boolean;
    terminalFontSize?: string | 'small' | 'medium' | 'big';
    terminalFontFamily?: string;
    terminalLineHeight?: number;
    logger?: (...params: any) => void;
  }

  interface Player {
    create(src: string, target?: HTMLElement, options?: PlayerCreateOptions): Player;
    dispose(): void;
  }

  export function create(src: string, target?: HTMLElement, options?: PlayerCreateOptions): Player;
  export function dispose(): void;
}
