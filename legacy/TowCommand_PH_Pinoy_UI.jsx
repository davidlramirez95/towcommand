import { useState, useEffect } from "react";

// ══════════════════════════════════════════════════════════════
// TOWCOMMAND PH — PINOY EDITION UI/UX
// ══════════════════════════════════════════════════════════════

const S = {
  LOGO:"logo",SPLASH:"splash",LOGIN:"login",HOME:"home",DIAGNOSE:"diagnose",
  SERVICE:"service",VEHICLE:"vehicle",DROPOFF:"dropoff",PRICE:"price",
  MATCHING:"matching",MATCHED:"matched",TRACKING:"tracking",CHAT:"chat",
  CONDITION:"condition",COMPLETE:"complete",RATE:"rate",SOS:"sos",
  PROVIDERS:"providers",HISTORY:"history",PROFILE:"profile",TYPHOON:"typhoon",SUKI:"suki",
};

const C = {
  navy:"#0B1D33",teal:"#00897B",gold:"#F5A623",orange:"#FF6B35",coral:"#FF4757",
  white:"#FFFFFF",cream:"#FFF9F0",light:"#F5F1EB",grey:"#8E8E93",greyL:"#E8E4DE",
  dark:"#1A1A2E",green:"#00C48C",blue:"#2196F3",bg:"#FAF7F2",
  g1:"linear-gradient(135deg,#0B1D33 0%,#1B3A5C 50%,#00897B 100%)",
  g2:"linear-gradient(135deg,#FF6B35 0%,#F5A623 100%)",
  g3:"linear-gradient(135deg,#00897B 0%,#00BFA5 100%)",
};

const F = { d:"'Poppins',sans-serif", b:"'Poppins',sans-serif", m:"'SF Mono',monospace" };

// ── SVG LOGO ──
const Logo = ({size=80,variant="full",anim=false}) => {
  const is = variant==="icon"?size:size*0.6;
  const icon = (
    <svg width={is} height={is} viewBox="0 0 120 120" fill="none">
      <defs>
        <linearGradient id="lg1" x1="0" y1="0" x2="120" y2="120" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#FF6B35"/><stop offset="50%" stopColor="#F5A623"/><stop offset="100%" stopColor="#FFD93D"/>
        </linearGradient>
        <linearGradient id="lg2" x1="0" y1="0" x2="120" y2="120" gradientUnits="userSpaceOnUse">
          <stop offset="0%" stopColor="#00897B"/><stop offset="100%" stopColor="#00BFA5"/>
        </linearGradient>
        <filter id="ls" x="-10%" y="-10%" width="130%" height="130%">
          <feDropShadow dx="0" dy="4" stdDeviation="8" floodColor="#FF6B35" floodOpacity="0.3"/>
        </filter>
      </defs>
      <circle cx="60" cy="60" r="56" fill="#0B1D33"/>
      <circle cx="60" cy="60" r="52" fill="url(#lg2)" opacity="0.15"/>
      <g filter="url(#ls)">
        <path d="M45 35L45 65Q45 82 60 82Q75 82 75 65L75 58" stroke="url(#lg1)" strokeWidth="8" strokeLinecap="round" fill="none"/>
        <path d="M35 35L55 35" stroke="url(#lg1)" strokeWidth="8" strokeLinecap="round"/>
        <path d="M75 55L75 30M65 40L75 30L85 40" stroke="#FFD93D" strokeWidth="6" strokeLinecap="round" strokeLinejoin="round"/>
      </g>
      {[0,1,2].map(i=><line key={i} x1={82+i*6} y1={88} x2={85+i*8} y2={95+i*2} stroke="#F5A623" strokeWidth="2.5" strokeLinecap="round" opacity={0.6-i*0.15}/>)}
      <g transform="translate(60,98) scale(0.5)">
        <polygon points="0,-10 2.9,-4 9.5,-3.1 4.8,1.5 5.9,8.1 0,5 -5.9,8.1 -4.8,1.5 -9.5,-3.1 -2.9,-4" fill="#F5A623" opacity="0.8"/>
      </g>
    </svg>
  );
  if(variant==="icon") return icon;
  return (
    <div style={{display:"flex",flexDirection:"column",alignItems:"center",gap:variant==="splash"?16:8}}>
      <div style={{animation:anim?"logoPulse 2s ease-in-out infinite":"none",filter:anim?"drop-shadow(0 0 20px rgba(255,107,53,0.4))":"none"}}>{icon}</div>
      <div style={{textAlign:"center"}}>
        <div style={{fontFamily:F.d,fontSize:variant==="splash"?28:18,fontWeight:800,letterSpacing:-0.5,background:C.g2,WebkitBackgroundClip:"text",WebkitTextFillColor:"transparent"}}>TowCommand</div>
        <div style={{fontFamily:F.d,fontSize:variant==="splash"?13:9,fontWeight:600,color:C.teal,letterSpacing:4,marginTop:-2}}>PILIPINAS</div>
      </div>
    </div>
  );
};

// ── Shared Components ──
const Bar = ({light=false}) => (
  <div style={{display:"flex",justifyContent:"space-between",padding:"8px 20px 4px",fontSize:12,fontWeight:600,color:light?"#fff":C.navy,fontFamily:F.b}}>
    <span>9:41</span>
    <div style={{display:"flex",gap:5,alignItems:"center"}}>
      <svg width="17" height="12" viewBox="0 0 17 12"><rect x="0" y="7" width="3" height="5" rx="0.5" fill={light?"#fff":C.navy} opacity="0.4"/><rect x="4.5" y="4.5" width="3" height="7.5" rx="0.5" fill={light?"#fff":C.navy} opacity="0.6"/><rect x="9" y="2" width="3" height="10" rx="0.5" fill={light?"#fff":C.navy} opacity="0.8"/><rect x="13.5" y="0" width="3" height="12" rx="0.5" fill={light?"#fff":C.navy}/></svg>
      <div style={{width:25,height:12,border:`2px solid ${light?"rgba(255,255,255,0.6)":C.navy}`,borderRadius:4,position:"relative"}}>
        <div style={{position:"absolute",left:2,top:2,width:16,height:6,background:C.green,borderRadius:1.5}}/>
      </div>
    </div>
  </div>
);

const HI = () => <div style={{display:"flex",justifyContent:"center",padding:"8px 0 6px"}}><div style={{width:134,height:5,background:"rgba(0,0,0,0.15)",borderRadius:3}}/></div>;

const Btn = ({children,onClick,v="primary",full=false,sm=false,st={}}) => {
  const styles = {
    primary:{background:C.g2,color:"#fff",boxShadow:"0 4px 15px rgba(255,107,53,0.35)"},
    secondary:{background:C.light,color:C.navy,border:`1.5px solid ${C.greyL}`},
    teal:{background:C.g3,color:"#fff",boxShadow:"0 4px 15px rgba(0,137,123,0.3)"},
    danger:{background:C.coral,color:"#fff",boxShadow:"0 4px 15px rgba(255,71,87,0.35)"},
    ghost:{background:"transparent",color:C.orange},
  };
  return <button onClick={onClick} style={{fontFamily:F.d,fontWeight:700,border:"none",cursor:"pointer",display:"flex",alignItems:"center",justifyContent:"center",gap:8,borderRadius:sm?10:14,padding:sm?"8px 16px":"14px 24px",fontSize:sm?12:14,width:full?"100%":"auto",transition:"all .2s",...(styles[v]||{}),...st}}>{children}</button>;
};

const Card = ({children,onClick,st={},elevated=false,selected=false}) => (
  <div onClick={onClick} style={{background:C.white,borderRadius:16,padding:"14px 16px",border:selected?`2px solid ${C.orange}`:`1px solid ${C.greyL}`,boxShadow:elevated?"0 8px 30px rgba(11,29,51,0.08)":"0 2px 8px rgba(11,29,51,0.04)",cursor:onClick?"pointer":"default",transition:"all .2s",...(selected&&{background:C.cream}),...st}}>{children}</div>
);

const Nav = ({active="home",go}) => {
  const items = [{id:"home",icon:"🏠",lb:"Home"},{id:"history",icon:"📋",lb:"Activity"},{id:"suki",icon:"⭐",lb:"Suki"},{id:"profile",icon:"👤",lb:"Account"}];
  return (
    <div style={{display:"flex",borderTop:`1px solid ${C.greyL}`,background:C.white,padding:"4px 0 2px"}}>
      {items.map(it=>{const a=active===it.id;return(
        <div key={it.id} onClick={()=>go(S[it.id.toUpperCase()])} style={{flex:1,display:"flex",flexDirection:"column",alignItems:"center",gap:2,padding:"8px 4px 4px",cursor:"pointer"}}>
          <div style={{fontSize:20,filter:a?"none":"grayscale(100%) opacity(0.4)",transform:a?"scale(1.1)":"scale(1)"}}>{it.icon}</div>
          <span style={{fontFamily:F.b,fontSize:10,fontWeight:a?700:500,color:a?C.orange:C.grey}}>{it.lb}</span>
          {a&&<div style={{width:4,height:4,borderRadius:2,background:C.orange}}/>}
        </div>
      );})}
    </div>
  );
};

const Back = ({title,onBack,right}) => (
  <div style={{display:"flex",alignItems:"center",padding:"4px 16px 8px",gap:10}}>
    <div onClick={onBack} style={{width:36,height:36,borderRadius:12,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",cursor:"pointer",fontSize:16}}>←</div>
    <span style={{flex:1,fontFamily:F.d,fontSize:16,fontWeight:700,color:C.navy}}>{title}</span>
    {right}
  </div>
);

const AiBadge = () => <span style={{fontFamily:F.b,fontSize:8,fontWeight:700,color:"#fff",background:"linear-gradient(135deg,#667eea,#764ba2)",borderRadius:6,padding:"3px 8px"}}>AI-POWERED</span>;
// ── SCREEN: LOGO SHOWCASE ──
const LogoShowcase = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg,overflow:"auto"}}>
    <Bar/>
    <div style={{padding:"12px 20px",textAlign:"center"}}>
      <div style={{fontFamily:F.d,fontSize:18,fontWeight:800,color:C.navy,marginBottom:4}}>TowCommand PH</div>
      <div style={{fontFamily:F.b,fontSize:11,color:C.grey}}>Official Logo & Brand Identity</div>
    </div>
    <div style={{margin:"0 16px 12px",background:C.navy,borderRadius:20,padding:"30px 20px",display:"flex",flexDirection:"column",alignItems:"center",gap:16}}>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:600,color:"rgba(255,255,255,0.4)",letterSpacing:2}}>PRIMARY — DARK BACKGROUND</div>
      <Logo size={100} variant="full" anim/>
    </div>
    <div style={{margin:"0 16px 12px",background:C.white,borderRadius:20,padding:"24px 20px",border:`1px solid ${C.greyL}`,display:"flex",flexDirection:"column",alignItems:"center",gap:12}}>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:600,color:C.grey,letterSpacing:2}}>ON LIGHT BACKGROUND</div>
      <Logo size={80} variant="full"/>
    </div>
    <div style={{margin:"0 16px 12px",display:"flex",gap:10}}>
      {[{bg:C.navy,lb:"APP ICON",s:52},{bg:C.white,lb:"FAVICON",s:40,border:true},{bg:C.g2,lb:"SPLASH",s:40}].map((v,i)=>(
        <div key={i} style={{flex:1,background:v.bg,borderRadius:16,padding:20,display:"flex",flexDirection:"column",alignItems:"center",gap:8,...(v.border?{border:`1px solid ${C.greyL}`}:{})}}>
          <div style={{fontFamily:F.b,fontSize:8,fontWeight:600,color:v.border?C.grey:"rgba(255,255,255,0.5)",letterSpacing:1.5}}>{v.lb}</div>
          <Logo size={v.s} variant="icon"/>
        </div>
      ))}
    </div>
    <div style={{margin:"0 16px 12px",background:C.white,borderRadius:16,padding:16,border:`1px solid ${C.greyL}`}}>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:600,color:C.grey,letterSpacing:2,marginBottom:10}}>BRAND PALETTE</div>
      <div style={{display:"flex",gap:6,flexWrap:"wrap"}}>
        {[{c:C.navy,n:"Navy",h:"#0B1D33"},{c:C.teal,n:"Teal",h:"#00897B"},{c:C.orange,n:"Orange",h:"#FF6B35"},{c:C.gold,n:"Gold",h:"#F5A623"},{c:C.coral,n:"Coral",h:"#FF4757"},{c:C.green,n:"Success",h:"#00C48C"}].map(cl=>(
          <div key={cl.n} style={{flex:"1 0 28%",display:"flex",alignItems:"center",gap:6,padding:6,borderRadius:8,background:C.light}}>
            <div style={{width:24,height:24,borderRadius:6,background:cl.c,flexShrink:0}}/>
            <div><div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.navy}}>{cl.n}</div><div style={{fontFamily:F.m,fontSize:8,color:C.grey}}>{cl.h}</div></div>
          </div>
        ))}
      </div>
    </div>
    <div style={{padding:"0 16px 16px"}}><Btn full onClick={()=>go(S.SPLASH)}>Continue to App →</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: SPLASH ──
const Splash = ({go}) => {
  useEffect(()=>{const t=setTimeout(()=>go(S.LOGIN),2200);return()=>clearTimeout(t);},[]);
  return (
    <div style={{flex:1,display:"flex",flexDirection:"column",alignItems:"center",justifyContent:"center",background:C.g1,position:"relative",overflow:"hidden"}}>
      <div style={{position:"absolute",inset:0,opacity:0.05,background:"repeating-linear-gradient(45deg,transparent,transparent 35px,rgba(255,255,255,0.1) 35px,rgba(255,255,255,0.1) 70px)"}}/>
      <div style={{position:"absolute",top:-60,right:-60,width:200,height:200,borderRadius:100,background:"radial-gradient(circle,rgba(245,166,35,0.15) 0%,transparent 70%)"}}/>
      <Logo size={120} variant="splash" anim/>
      <div style={{position:"absolute",bottom:60,display:"flex",gap:6}}>
        {[0,1,2].map(i=><div key={i} style={{width:8,height:8,borderRadius:4,background:C.gold,opacity:0.3+i*0.3}}/>)}
      </div>
      <div style={{position:"absolute",bottom:30,fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.3)",letterSpacing:1}}>Tulong sa daan, anytime.</div>
      <HI/>
    </div>
  );
};

// ── SCREEN: LOGIN (Cognito SSO) ──
const Login = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{flex:1,display:"flex",flexDirection:"column",alignItems:"center",justifyContent:"center",padding:"0 28px"}}>
      <Logo size={90} variant="full"/>
      <div style={{marginTop:28,marginBottom:32,textAlign:"center"}}>
        <div style={{fontFamily:F.d,fontSize:20,fontWeight:800,color:C.navy,lineHeight:1.3}}>Get help on the road,<br/>anytime, anywhere.</div>
        <div style={{fontFamily:F.b,fontSize:12,color:C.grey,marginTop:8}}>Sign in to book a tow truck in seconds</div>
      </div>
      <div style={{width:"100%",display:"flex",flexDirection:"column",gap:10}}>
        <button onClick={()=>go(S.HOME)} style={{width:"100%",padding:14,borderRadius:14,border:`1.5px solid ${C.greyL}`,background:C.white,cursor:"pointer",display:"flex",alignItems:"center",justifyContent:"center",gap:10,fontFamily:F.d,fontSize:14,fontWeight:600,color:C.navy,boxShadow:"0 2px 8px rgba(0,0,0,0.04)"}}>
          <svg width="20" height="20" viewBox="0 0 24 24"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
          Continue with Google
        </button>
        <button onClick={()=>go(S.HOME)} style={{width:"100%",padding:14,borderRadius:14,border:"none",background:"#1877F2",cursor:"pointer",display:"flex",alignItems:"center",justifyContent:"center",gap:10,fontFamily:F.d,fontSize:14,fontWeight:600,color:"#fff",boxShadow:"0 4px 15px rgba(24,119,242,0.3)"}}>
          <svg width="20" height="20" viewBox="0 0 24 24" fill="white"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>
          Continue with Facebook
        </button>
        <button onClick={()=>go(S.HOME)} style={{width:"100%",padding:14,borderRadius:14,border:"none",background:C.navy,cursor:"pointer",display:"flex",alignItems:"center",justifyContent:"center",gap:10,fontFamily:F.d,fontSize:14,fontWeight:600,color:"#fff"}}>
          <svg width="18" height="22" viewBox="0 0 18 22" fill="white"><path d="M14.94 11.58c-.03-2.73 2.22-4.04 2.32-4.1-1.27-1.85-3.24-2.1-3.94-2.13-1.67-.17-3.28 1-4.13 1-.86 0-2.17-.97-3.57-.95-1.83.03-3.53 1.07-4.47 2.72-1.91 3.32-.49 8.24 1.37 10.93.91 1.32 2 2.8 3.43 2.75 1.38-.06 1.9-.89 3.56-.89s2.13.89 3.58.86c1.48-.03 2.39-1.34 3.29-2.67 1.04-1.53 1.46-3.01 1.49-3.09-.03-.01-2.86-1.1-2.89-4.35z"/></svg>
          Continue with Apple
        </button>
      </div>
      <div style={{display:"flex",alignItems:"center",gap:12,width:"100%",margin:"20px 0"}}>
        <div style={{flex:1,height:1,background:C.greyL}}/><span style={{fontFamily:F.b,fontSize:11,color:C.grey}}>or</span><div style={{flex:1,height:1,background:C.greyL}}/>
      </div>
      <div style={{width:"100%",display:"flex",gap:8}}>
        <div style={{width:70,padding:"12px 8px",borderRadius:12,border:`1.5px solid ${C.greyL}`,background:C.white,display:"flex",alignItems:"center",justifyContent:"center",gap:4,fontFamily:F.b,fontSize:13,fontWeight:600,color:C.navy}}>🇵🇭 +63</div>
        <div onClick={()=>go(S.HOME)} style={{flex:1,padding:"12px 14px",borderRadius:12,border:`1.5px solid ${C.greyL}`,background:C.white,fontFamily:F.b,fontSize:13,color:C.grey,cursor:"pointer",display:"flex",alignItems:"center"}}>9XX XXX XXXX</div>
      </div>
    </div>
    <div style={{padding:"0 28px 8px",textAlign:"center"}}>
      <span style={{fontFamily:F.b,fontSize:10,color:C.grey}}>By signing in, you agree to our <span style={{color:C.orange,fontWeight:600}}>Terms</span> and <span style={{color:C.orange,fontWeight:600}}>Privacy Policy</span></span>
    </div>
    <HI/>
  </div>
);

// ── SCREEN: HOME ──
const Home = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"4px 20px 0",flex:1,overflow:"auto"}}>
      <div style={{display:"flex",justifyContent:"space-between",alignItems:"center",marginBottom:14}}>
        <div>
          <div style={{fontFamily:F.b,fontSize:11,color:C.grey}}>Magandang hapon 🌤️</div>
          <div style={{fontFamily:F.d,fontSize:18,fontWeight:800,color:C.navy}}>David!</div>
        </div>
        <div style={{display:"flex",gap:8}}>
          <div style={{width:36,height:36,borderRadius:12,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",fontSize:16}}>🔔</div>
          <div style={{width:36,height:36,borderRadius:12,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:14,color:"#fff",fontWeight:800,fontFamily:F.d}}>D</div>
        </div>
      </div>
      
      <Card onClick={()=>go(S.DIAGNOSE)} elevated st={{background:C.g1,padding:18,marginBottom:14,border:"none",position:"relative",overflow:"hidden"}}>
        <div style={{position:"absolute",top:-20,right:-20,width:100,height:100,borderRadius:50,background:"radial-gradient(circle,rgba(245,166,35,0.2) 0%,transparent 70%)"}}/>
        <div style={{display:"flex",alignItems:"center",gap:12}}>
          <div style={{width:48,height:48,borderRadius:14,background:"rgba(255,107,53,0.2)",display:"flex",alignItems:"center",justifyContent:"center",fontSize:26,flexShrink:0}}>🤖</div>
          <div style={{flex:1}}>
            <div style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:"#fff"}}>What's wrong with your car?</div>
            <div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.6)",marginTop:2}}>AI will diagnose & find the cheapest fix</div>
          </div>
          <div style={{fontSize:20,color:C.gold}}>→</div>
        </div>
      </Card>

      <div style={{height:145,borderRadius:18,overflow:"hidden",marginBottom:14,position:"relative",background:"linear-gradient(135deg,#E8E4DE,#D4CEC6)"}}>
        {[0,1,2,3,4,5].map(i=><div key={i} style={{position:"absolute",left:0,right:0,top:i*29,height:1,background:"rgba(11,29,51,0.06)"}}/>)}
        {[0,1,2,3,4,5,6,7].map(i=><div key={i} style={{position:"absolute",top:0,bottom:0,left:i*45,width:1,background:"rgba(11,29,51,0.06)"}}/>)}
        <div style={{position:"absolute",top:50,left:0,right:0,height:2,background:"rgba(11,29,51,0.12)"}}/>
        <div style={{position:"absolute",top:0,bottom:0,left:"55%",width:2,background:"rgba(11,29,51,0.12)"}}/>
        <div style={{position:"absolute",top:44,left:"53%",transform:"translateX(-50%)"}}>
          <div style={{width:24,height:24,borderRadius:12,background:C.orange,border:"3px solid #fff",boxShadow:"0 3px 10px rgba(255,107,53,0.4)"}}><div style={{width:6,height:6,borderRadius:3,background:"#fff",margin:"6px auto"}}/></div>
        </div>
        <div style={{position:"absolute",top:25,left:"25%",fontSize:18,filter:"drop-shadow(0 2px 4px rgba(0,0,0,0.2))"}}>🚛</div>
        <div style={{position:"absolute",top:80,left:"70%",fontSize:16,filter:"drop-shadow(0 2px 4px rgba(0,0,0,0.2))"}}>🚛</div>
        <div style={{position:"absolute",top:100,left:"35%",fontSize:14,opacity:0.6}}>🚛</div>
        <div style={{position:"absolute",top:10,right:10,background:"rgba(11,29,51,0.85)",backdropFilter:"blur(10px)",borderRadius:10,padding:"6px 10px",display:"flex",alignItems:"center",gap:5}}>
          <div style={{width:6,height:6,borderRadius:3,background:C.green}}/>
          <span style={{fontFamily:F.b,fontSize:10,fontWeight:600,color:"#fff"}}>3 trucks nearby</span>
        </div>
      </div>

      <div style={{display:"flex",justifyContent:"space-between",alignItems:"center",marginBottom:10}}>
        <span style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:C.navy}}>Quick Services</span>
        <span style={{fontFamily:F.b,fontSize:11,fontWeight:600,color:C.orange}}>See all</span>
      </div>
      <div style={{display:"grid",gridTemplateColumns:"1fr 1fr 1fr 1fr",gap:8,marginBottom:16}}>
        {[{i:"🚛",l:"Tow",c:"#FFF3EB"},{i:"⛽",l:"Fuel",c:"#E8F8F0"},{i:"🔋",l:"Jumpstart",c:"#EBF0FF"},{i:"🔧",l:"Mechanic",c:"#FFF8E1"}].map(s=>(
          <div key={s.l} onClick={()=>go(S.SERVICE)} style={{display:"flex",flexDirection:"column",alignItems:"center",gap:6,padding:"14px 8px",borderRadius:14,background:s.c,cursor:"pointer"}}>
            <span style={{fontSize:24}}>{s.i}</span>
            <span style={{fontFamily:F.b,fontSize:10,fontWeight:600,color:C.navy}}>{s.l}</span>
          </div>
        ))}
      </div>
      
      <Card st={{background:C.cream,marginBottom:12,border:`1.5px solid ${C.gold}33`}} onClick={()=>go(S.SUKI)}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <div style={{fontSize:22}}>⭐</div>
          <div style={{flex:1}}>
            <div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.navy}}>Suki Silver Member</div>
            <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>2 more bookings for Gold • 5% off all services</div>
          </div>
        </div>
      </Card>
    </div>
    <Nav active="home" go={go}/>
    <HI/>
  </div>
);

// ── SCREEN: SMART DIAGNOSIS ──
const Diagnose = ({go}) => {
  const [step,setStep]=useState(0);
  const [sel,setSel]=useState([]);
  const syms=[{ic:"⛽",lb:"Empty fuel",id:"fuel"},{ic:"💨",lb:"Flat tire",id:"tire"},{ic:"🔋",lb:"Dead battery",id:"battery"},{ic:"🔑",lb:"Locked out",id:"lockout"},{ic:"🌡️",lb:"Overheating",id:"overheat"},{ic:"🚨",lb:"Accident",id:"accident"},{ic:"⚙️",lb:"Engine trouble",id:"engine"},{ic:"💧",lb:"Leaking fluid",id:"leak"},{ic:"🔌",lb:"Electrical",id:"electrical"},{ic:"❓",lb:"Not sure",id:"unknown"}];
  const tog=id=>setSel(p=>p.includes(id)?p.filter(x=>x!==id):[...p,id]);
  useEffect(()=>{if(step===1){const t=setTimeout(()=>setStep(2),2000);return()=>clearTimeout(t);}},[step]);

  if(step===1) return (
    <div style={{flex:1,display:"flex",flexDirection:"column",background:C.g1,alignItems:"center",justifyContent:"center",gap:16}}>
      <div style={{position:"relative",width:90,height:90}}>
        <div style={{position:"absolute",inset:0,border:"3px solid rgba(245,166,35,0.2)",borderRadius:45}}/>
        <div style={{position:"absolute",inset:14,background:C.g2,borderRadius:35,display:"flex",alignItems:"center",justifyContent:"center",fontSize:34,boxShadow:"0 8px 30px rgba(255,107,53,0.4)"}}>🤖</div>
      </div>
      <div style={{fontFamily:F.d,fontSize:16,fontWeight:700,color:"#fff"}}>Analyzing your symptoms...</div>
      <div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.5)",textAlign:"center",padding:"0 50px"}}>AI is matching {sel.length} symptom{sel.length>1?"s":""} to the best service</div>
    </div>
  );

  if(step===2) return (
    <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
      <Bar/><Back title="AI Diagnosis Result" onBack={()=>setStep(0)} right={<AiBadge/>}/>
      <div style={{flex:1,overflow:"auto",padding:"0 16px 16px"}}>
        <Card elevated st={{background:"linear-gradient(135deg,#E8F8F0,#d4f5e2)",marginBottom:12,border:"1.5px solid rgba(0,196,140,0.3)"}}>
          <div style={{display:"flex",alignItems:"center",gap:10}}>
            <div style={{width:44,height:44,borderRadius:14,background:"rgba(0,196,140,0.15)",display:"flex",alignItems:"center",justifyContent:"center",fontSize:24}}>✅</div>
            <div><div style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:"#1A7F5A"}}>You DON'T need a tow!</div><div style={{fontFamily:F.b,fontSize:10,color:C.green}}>92% confidence • Save up to ₱2,500</div><div style={{fontFamily:F.b,fontSize:10,color:"#2D6A4F",marginTop:4}}>A roadside mechanic can fix this on-the-spot.</div></div>
          </div>
        </Card>
        <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.orange,letterSpacing:1.5,marginBottom:8}}>✨ RECOMMENDED FOR YOU</div>
        <Card onClick={()=>go(S.VEHICLE)} elevated selected st={{marginBottom:12,position:"relative"}}>
          <div style={{position:"absolute",top:10,right:10,fontFamily:F.b,fontSize:8,fontWeight:700,color:"#fff",background:C.green,borderRadius:6,padding:"3px 8px"}}>BEST MATCH</div>
          <div style={{display:"flex",alignItems:"center",gap:12}}>
            <div style={{width:46,height:46,borderRadius:14,background:"#FFF3EB",display:"flex",alignItems:"center",justifyContent:"center",fontSize:22}}>🔧</div>
            <div><div style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:C.navy}}>On-Site Mechanic</div><div style={{fontFamily:F.b,fontSize:11,color:C.grey,marginTop:2}}>Tire change + basic repair</div><div style={{display:"flex",gap:10,marginTop:5}}><span style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.green}}>₱500–₱1,200</span><span style={{fontFamily:F.b,fontSize:11,color:C.orange,fontWeight:600}}>ETA 15 min</span></div></div>
          </div>
        </Card>
        <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>OTHER OPTIONS</div>
        {[{ic:"⛽",t:"Fuel Delivery",d:"Gasolina delivered to you",p:"₱200–₱500",eta:"20 min"},{ic:"🚛",t:"Full Tow",d:"If repair doesn't work",p:"₱1,800–₱3,500",eta:"12 min"}].map((o,i)=>(
          <Card key={i} onClick={()=>go(S.VEHICLE)} st={{marginBottom:8}}>
            <div style={{display:"flex",alignItems:"center",gap:10}}>
              <div style={{width:40,height:40,borderRadius:12,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",fontSize:18}}>{o.ic}</div>
              <div style={{flex:1}}><div style={{fontFamily:F.d,fontSize:12,fontWeight:600,color:C.navy}}>{o.t}</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>{o.d}</div></div>
              <div style={{textAlign:"right"}}><div style={{fontFamily:F.d,fontSize:11,fontWeight:700,color:C.navy}}>{o.p}</div><div style={{fontFamily:F.b,fontSize:9,color:C.grey}}>{o.eta}</div></div>
            </div>
          </Card>
        ))}
        <Card st={{background:C.cream,marginTop:4,border:`1.5px solid ${C.gold}33`}}>
          <div style={{display:"flex",alignItems:"center",gap:10}}><span style={{fontSize:20}}>💰</span><div><div style={{fontFamily:F.d,fontSize:11,fontWeight:700,color:"#F57F17"}}>Smart Savings</div><div style={{fontFamily:F.b,fontSize:10,color:"#795548"}}>Users who used AI Diagnosis saved avg <b>₱2,100</b> per incident</div></div></div>
        </Card>
        <div style={{textAlign:"center",marginTop:12}}><span onClick={()=>go(S.SERVICE)} style={{fontFamily:F.b,fontSize:11,color:C.grey,textDecoration:"underline",cursor:"pointer"}}>Skip AI → Choose service manually</span></div>
      </div><HI/>
    </div>
  );

  return (
    <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
      <Bar/><Back title="Smart Diagnosis" onBack={()=>go(S.HOME)} right={<AiBadge/>}/>
      <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
        <Card elevated st={{background:C.g1,marginBottom:14,border:"none"}}>
          <div style={{display:"flex",alignItems:"center",gap:10,marginBottom:8}}><span style={{fontSize:22}}>🤖</span><div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:"#fff"}}>What's happening to your vehicle?</div></div>
          <div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.6)",lineHeight:1.5}}>Select all symptoms so I can recommend the <b style={{color:"#fff"}}>cheapest & fastest</b> solution.</div>
        </Card>
        <div style={{display:"flex",gap:8,marginBottom:14}}>
          {[{i:"🎤",t:"Voice Describe",s:"Tagalog or English"},{i:"📷",t:"Take Photo",s:"AI visual analysis"}].map(o=>(
            <Card key={o.t} st={{flex:1,textAlign:"center",padding:"14px 10px"}}><span style={{fontSize:22}}>{o.i}</span><div style={{fontFamily:F.d,fontSize:10,fontWeight:600,color:C.navy,marginTop:4}}>{o.t}</div><div style={{fontFamily:F.b,fontSize:9,color:C.grey}}>{o.s}</div></Card>
          ))}
        </div>
        <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:10}}>OR TAP YOUR SYMPTOMS</div>
        <div style={{display:"flex",flexWrap:"wrap",gap:7}}>
          {syms.map(s=>(
            <div key={s.id} onClick={()=>tog(s.id)} style={{display:"flex",alignItems:"center",gap:7,padding:"10px 13px",borderRadius:12,border:sel.includes(s.id)?`2px solid ${C.orange}`:`1.5px solid ${C.greyL}`,background:sel.includes(s.id)?C.cream:C.white,cursor:"pointer",boxShadow:sel.includes(s.id)?"0 2px 10px rgba(255,107,53,0.15)":"none"}}>
              <span style={{fontSize:17}}>{s.ic}</span><span style={{fontFamily:F.b,fontSize:11,fontWeight:sel.includes(s.id)?700:500,color:sel.includes(s.id)?C.orange:C.navy}}>{s.lb}</span>
            </div>
          ))}
        </div>
        {sel.length>0&&<Card st={{marginTop:12,background:C.cream,border:`1px solid ${C.gold}33`}}><div style={{display:"flex",alignItems:"center",gap:8}}><span style={{fontSize:16}}>💡</span><div style={{fontFamily:F.b,fontSize:10.5,color:"#5D6D7E"}}>{sel.includes("fuel")?"This looks like you just need fuel delivery!":sel.includes("tire")?"Flat tire? On-site mechanic can change this in 15 mins.":sel.includes("accident")?"Accident detected. We'll prioritize emergency recovery.":"Analyzing your combination of symptoms..."}</div></div></Card>}
      </div>
      <div style={{padding:"12px 16px 4px"}}>
        <Btn full onClick={()=>sel.length>0&&setStep(1)} st={{opacity:sel.length>0?1:0.4}}>🤖 Diagnose My Problem ({sel.length})</Btn>
        <div style={{textAlign:"center",marginTop:10}}><span onClick={()=>go(S.SERVICE)} style={{fontFamily:F.b,fontSize:11,color:C.grey,textDecoration:"underline",cursor:"pointer"}}>I already know what I need →</span></div>
      </div><HI/>
    </div>
  );
};
// ── SCREEN: SERVICE SELECTION ──
const Service = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Choose Service" onBack={()=>go(S.HOME)}/>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      {[
        {ic:"🚛",t:"Flatbed Tow",d:"Full vehicle transport, safest option",p:"₱1,800–₱3,500",eta:"10–15 min",pop:true},
        {ic:"🔗",t:"Wheel-Lift Tow",d:"2-wheel lift, budget-friendly",p:"₱1,200–₱2,200",eta:"8–12 min"},
        {ic:"⛽",t:"Fuel Delivery",d:"Gas/diesel delivered to your location",p:"₱200–₱500",eta:"15–20 min"},
        {ic:"🔋",t:"Jumpstart",d:"Dead battery? Get going in minutes",p:"₱300–₱600",eta:"10–15 min"},
        {ic:"🛞",t:"Tire Change",d:"Flat tire replacement on-site",p:"₱400–₱800",eta:"15–20 min"},
        {ic:"🔑",t:"Lockout Service",d:"Locked out? We'll get you in",p:"₱500–₱1,000",eta:"15–25 min"},
      ].map((s,i)=>(
        <Card key={i} onClick={()=>go(S.VEHICLE)} st={{marginBottom:8,position:"relative"}}>
          {s.pop&&<div style={{position:"absolute",top:10,right:10,fontFamily:F.b,fontSize:8,fontWeight:700,color:"#fff",background:C.orange,borderRadius:6,padding:"2px 8px"}}>POPULAR</div>}
          <div style={{display:"flex",alignItems:"center",gap:12}}>
            <div style={{width:46,height:46,borderRadius:14,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",fontSize:22,flexShrink:0}}>{s.ic}</div>
            <div style={{flex:1}}>
              <div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy}}>{s.t}</div>
              <div style={{fontFamily:F.b,fontSize:10,color:C.grey,marginTop:1}}>{s.d}</div>
              <div style={{display:"flex",gap:10,marginTop:4}}>
                <span style={{fontFamily:F.d,fontSize:11,fontWeight:700,color:C.green}}>{s.p}</span>
                <span style={{fontFamily:F.b,fontSize:10,color:C.orange,fontWeight:600}}>ETA {s.eta}</span>
              </div>
            </div>
            <div style={{fontSize:16,color:C.greyL}}>→</div>
          </div>
        </Card>
      ))}
    </div><HI/>
  </div>
);

// ── SCREEN: VEHICLE ──
const Vehicle = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Vehicle Details" onBack={()=>go(S.SERVICE)}/>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>SAVED VEHICLES</div>
      <Card selected st={{marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:12}}>
          <div style={{width:46,height:46,borderRadius:14,background:C.cream,display:"flex",alignItems:"center",justifyContent:"center",fontSize:22}}>🚗</div>
          <div style={{flex:1}}>
            <div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy}}>2026 Montero GLS</div>
            <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>ABC 1234 • Midsize SUV • White</div>
          </div>
          <div style={{width:20,height:20,borderRadius:10,background:C.orange,display:"flex",alignItems:"center",justifyContent:"center"}}><span style={{color:"#fff",fontSize:12,fontWeight:800}}>✓</span></div>
        </div>
      </Card>
      <Card st={{marginBottom:14,border:`1.5px dashed ${C.greyL}`,textAlign:"center",padding:"16px"}}>
        <span style={{fontSize:22}}>➕</span>
        <div style={{fontFamily:F.d,fontSize:12,fontWeight:600,color:C.navy,marginTop:4}}>Add New Vehicle</div>
      </Card>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>VEHICLE CONDITION</div>
      {["Engine won't start","Flat tire(s)","Accident damage","Other / Not sure"].map((c,i)=>(
        <Card key={i} st={{marginBottom:6,padding:"12px 14px",...(i===0?{border:`2px solid ${C.orange}`,background:C.cream}:{})}}>
          <div style={{display:"flex",alignItems:"center",justifyContent:"space-between"}}>
            <span style={{fontFamily:F.b,fontSize:12,fontWeight:i===0?700:500,color:i===0?C.orange:C.navy}}>{c}</span>
            {i===0&&<div style={{width:18,height:18,borderRadius:9,background:C.orange,display:"flex",alignItems:"center",justifyContent:"center"}}><span style={{color:"#fff",fontSize:10}}>✓</span></div>}
          </div>
        </Card>
      ))}
    </div>
    <div style={{padding:"12px 16px 4px"}}><Btn full onClick={()=>go(S.DROPOFF)}>Continue →</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: DROPOFF ──
const Dropoff = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Drop-off Location" onBack={()=>go(S.VEHICLE)}/>
    <div style={{flex:1,padding:"0 16px"}}>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>PICKUP LOCATION</div>
      <Card st={{marginBottom:12}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <div style={{width:10,height:10,borderRadius:5,background:C.orange,flexShrink:0}}/>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.navy}}>Current Location</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>EDSA cor. Ayala Ave, Makati City</div></div>
        </div>
      </Card>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>DROP-OFF LOCATION</div>
      <Card selected st={{marginBottom:12}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <div style={{width:10,height:10,borderRadius:5,background:C.teal,flexShrink:0}}/>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.navy}}>Toyota Shaw, Mandaluyong</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>EDSA cor. Shaw Blvd • 4.2 km</div></div>
        </div>
      </Card>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>RECENT LOCATIONS</div>
      {["Casa — BF Homes, Parañaque","Mitsubishi Ortigas","AutoHub SLEX, Alabang"].map((l,i)=>(
        <Card key={i} st={{marginBottom:6,padding:"11px 14px"}}>
          <div style={{display:"flex",alignItems:"center",gap:10}}>
            <span style={{fontSize:14}}>📍</span>
            <span style={{fontFamily:F.b,fontSize:11,color:C.navy}}>{l}</span>
          </div>
        </Card>
      ))}
    </div>
    <div style={{padding:"12px 16px 4px"}}><Btn full onClick={()=>go(S.PRICE)}>Confirm Route →</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: PRICE ──
const Price = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Price Estimate" onBack={()=>go(S.DROPOFF)}/>
    <div style={{flex:1,padding:"0 16px"}}>
      <Card elevated st={{marginBottom:14}}>
        <div style={{textAlign:"center",marginBottom:12}}>
          <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>Estimated Total</div>
          <div style={{fontFamily:F.d,fontSize:36,fontWeight:800,color:C.orange,marginTop:4}}>₱1,850</div>
          <div style={{fontFamily:F.b,fontSize:10,color:C.grey,marginTop:2}}>MMDA Reg. 24-004 compliant pricing</div>
        </div>
        <div style={{height:1,background:C.greyL,margin:"10px 0"}}/>
        {[["Base fare","₱800"],["Distance (4.2 km × ₱100)","₱420"],["Weight surcharge (SUV)","₱250"],["Time surcharge (peak)","₱180"],["Platform fee","₱200"]].map(([l,v],i)=>(
          <div key={i} style={{display:"flex",justifyContent:"space-between",padding:"6px 0"}}>
            <span style={{fontFamily:F.b,fontSize:11,color:C.grey}}>{l}</span>
            <span style={{fontFamily:F.b,fontSize:11,fontWeight:600,color:C.navy}}>{v}</span>
          </div>
        ))}
      </Card>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>PAYMENT METHOD</div>
      <Card selected st={{marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <div style={{width:36,height:36,borderRadius:10,background:"#00B4D8",display:"flex",alignItems:"center",justifyContent:"center",fontSize:12,color:"#fff",fontWeight:800,fontFamily:F.d}}>G</div>
          <div style={{flex:1}}><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.navy}}>GCash</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>**** 8847</div></div>
          <div style={{width:18,height:18,borderRadius:9,background:C.orange,display:"flex",alignItems:"center",justifyContent:"center"}}><span style={{color:"#fff",fontSize:10}}>✓</span></div>
        </div>
      </Card>
      <Card st={{background:C.cream,border:`1.5px solid ${C.gold}33`}}>
        <div style={{display:"flex",alignItems:"center",gap:8}}>
          <span style={{fontSize:14}}>💡</span>
          <div style={{fontFamily:F.b,fontSize:10,color:"#5D6D7E"}}>A <b>₱200 hold</b> will be placed on your GCash. Final amount charged after job completion.</div>
        </div>
      </Card>
    </div>
    <div style={{padding:"12px 16px 4px"}}><Btn full onClick={()=>go(S.MATCHING)}>🚛 Book Now — ₱1,850</Btn></div>
    <HI/>
  </div>
);
// ── SCREEN: MATCHING ──
const Matching = ({go}) => {
  useEffect(()=>{const t=setTimeout(()=>go(S.MATCHED),3000);return()=>clearTimeout(t);},[]);
  return (
    <div style={{flex:1,display:"flex",flexDirection:"column",background:C.g1,alignItems:"center",justifyContent:"center",position:"relative"}}>
      <div style={{position:"absolute",inset:0,opacity:0.05,background:"repeating-linear-gradient(45deg,transparent,transparent 35px,rgba(255,255,255,0.1) 35px,rgba(255,255,255,0.1) 70px)"}}/>
      <div style={{position:"relative",width:100,height:100,marginBottom:20}}>
        <div style={{position:"absolute",inset:0,border:"3px solid rgba(245,166,35,0.15)",borderRadius:50}}/>
        <div style={{position:"absolute",inset:18,background:C.g2,borderRadius:35,display:"flex",alignItems:"center",justifyContent:"center",fontSize:38,boxShadow:"0 8px 30px rgba(255,107,53,0.4)"}}>🚛</div>
      </div>
      <div style={{fontFamily:F.d,fontSize:18,fontWeight:700,color:"#fff",marginBottom:6}}>Finding nearby trucks...</div>
      <div style={{fontFamily:F.b,fontSize:12,color:"rgba(255,255,255,0.5)",marginBottom:24}}>Matching you with the best provider</div>
      <div style={{display:"flex",gap:10}}>
        {["📍 Makati","🚛 3 available","⏱ ~2 min"].map(i=>(
          <div key={i} style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)",background:"rgba(255,255,255,0.08)",padding:"6px 12px",borderRadius:8}}>{i}</div>
        ))}
      </div>
      <div style={{position:"absolute",bottom:60}}><Btn v="ghost" onClick={()=>go(S.HOME)} sm st={{color:"rgba(255,255,255,0.5)"}}>Cancel Search</Btn></div>
      <HI/>
    </div>
  );
};

// ── SCREEN: MATCHED ──
const Matched = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"8px 16px",flex:1,overflow:"auto"}}>
      <div style={{height:155,borderRadius:18,overflow:"hidden",marginBottom:12,position:"relative",background:"linear-gradient(135deg,#E8E4DE,#D4CEC6)"}}>
        {[0,1,2,3,4,5].map(i=><div key={i} style={{position:"absolute",left:0,right:0,top:i*31,height:1,background:"rgba(11,29,51,0.06)"}}/>)}
        <div style={{position:"absolute",top:"40%",left:"55%",transform:"translate(-50%,-50%)",width:20,height:20,borderRadius:10,background:C.orange,border:"3px solid #fff",boxShadow:"0 3px 10px rgba(255,107,53,0.4)"}}/>
        <div style={{position:"absolute",top:"55%",left:"35%",fontSize:20}}>🚛</div>
        <div style={{position:"absolute",top:8,left:8,right:8,background:"rgba(11,29,51,0.85)",backdropFilter:"blur(10px)",borderRadius:12,padding:"10px 14px",display:"flex",alignItems:"center",gap:8}}>
          <div style={{width:8,height:8,borderRadius:4,background:C.green}}/>
          <span style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:"#fff"}}>Driver is on the way</span>
          <span style={{fontFamily:F.b,fontSize:12,fontWeight:700,color:C.gold,marginLeft:"auto"}}>ETA 8 min</span>
        </div>
      </div>
      <Card elevated st={{marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:12}}>
          <div style={{width:50,height:50,borderRadius:16,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:24,color:"#fff",fontWeight:800,fontFamily:F.d}}>JR</div>
          <div style={{flex:1}}>
            <div style={{display:"flex",alignItems:"center",gap:6}}>
              <span style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:C.navy}}>Juan Reyes</span>
              <span style={{fontFamily:F.b,fontSize:8,fontWeight:700,color:"#fff",background:C.teal,borderRadius:4,padding:"2px 6px"}}>VERIFIED</span>
            </div>
            <div style={{display:"flex",alignItems:"center",gap:8,marginTop:3}}>
              <span style={{fontFamily:F.b,fontSize:11,color:C.gold}}>★ 4.9</span>
              <span style={{fontFamily:F.b,fontSize:10,color:C.grey}}>• 847 jobs • ABC 1234</span>
            </div>
          </div>
        </div>
        <div style={{display:"flex",gap:8,marginTop:12}}>
          <Btn v="teal" sm full onClick={()=>go(S.CHAT)}>💬 Message</Btn>
          <Btn v="secondary" sm full onClick={()=>go(S.TRACKING)}>📍 Track</Btn>
        </div>
      </Card>
      <Card st={{background:C.cream,border:`1.5px solid ${C.gold}33`}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}><div style={{fontSize:20}}>🔐</div><div style={{flex:1}}><div style={{fontFamily:F.d,fontSize:11,fontWeight:700,color:C.navy}}>Digital Padala OTP</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>Share this code when driver arrives</div></div></div>
        <div style={{display:"flex",gap:6,marginTop:10,justifyContent:"center"}}>
          {"482917".split("").map((d,i)=>(
            <div key={i} style={{width:36,height:44,borderRadius:10,background:C.white,border:`2px solid ${C.orange}`,display:"flex",alignItems:"center",justifyContent:"center",fontFamily:F.d,fontSize:22,fontWeight:800,color:C.orange,boxShadow:"0 2px 8px rgba(255,107,53,0.15)"}}>{d}</div>
          ))}
        </div>
      </Card>
    </div>
    <div style={{padding:"8px 16px 4px"}}><Btn v="danger" full sm onClick={()=>go(S.SOS)}>🆘 Emergency SOS</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: TRACKING ──
const Tracking = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{flex:1,position:"relative",background:"linear-gradient(135deg,#E8E4DE,#D4CEC6)"}}>
      {[0,1,2,3,4,5,6,7,8].map(i=><div key={i} style={{position:"absolute",left:0,right:0,top:i*50,height:1,background:"rgba(11,29,51,0.06)"}}/>)}
      <svg style={{position:"absolute",inset:0,width:"100%",height:"100%"}} viewBox="0 0 400 400"><path d="M200 280Q200 200 160 160Q130 130 150 100" stroke={C.teal} strokeWidth="3" strokeDasharray="8 4" fill="none"/></svg>
      <div style={{position:"absolute",top:"65%",left:"48%",transform:"translate(-50%,-50%)"}}><div style={{width:24,height:24,borderRadius:12,background:C.orange,border:"3px solid #fff",boxShadow:"0 3px 12px rgba(255,107,53,0.4)"}}><div style={{width:6,height:6,borderRadius:3,background:"#fff",margin:"6px auto"}}/></div></div>
      <div style={{position:"absolute",top:"22%",left:"35%",transform:"translate(-50%,-50%)",fontSize:28,filter:"drop-shadow(0 3px 8px rgba(0,0,0,0.3))"}}>🚛</div>
      <div style={{position:"absolute",top:10,left:10,right:10,background:"rgba(11,29,51,0.9)",backdropFilter:"blur(12px)",borderRadius:14,padding:"12px 16px"}}>
        <div style={{display:"flex",justifyContent:"space-between",alignItems:"center"}}>
          <div><div style={{fontFamily:F.d,fontSize:14,fontWeight:700,color:"#fff"}}>Driver en route</div><div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.5)"}}>Juan Reyes • ABC 1234</div></div>
          <div style={{textAlign:"right"}}><div style={{fontFamily:F.d,fontSize:20,fontWeight:800,color:C.gold}}>5 min</div><div style={{fontFamily:F.b,fontSize:9,color:"rgba(255,255,255,0.4)"}}>ETA to you</div></div>
        </div>
        <div style={{display:"flex",gap:4,marginTop:10}}>
          {["Matched","En Route","Arrived","Loading","Complete"].map((s,i)=><div key={s} style={{flex:1,height:3,borderRadius:2,background:i<2?C.green:"rgba(255,255,255,0.15)"}}/>)}
        </div>
      </div>
      <div style={{position:"absolute",bottom:10,right:10}}><Btn v="secondary" sm>📤 Share Trip</Btn></div>
      <div style={{position:"absolute",bottom:10,left:10}}><Btn v="danger" sm onClick={()=>go(S.SOS)}>🆘</Btn></div>
    </div>
    <div style={{padding:"12px 16px 4px",display:"flex",gap:8}}>
      <Btn v="teal" full onClick={()=>go(S.CHAT)}>💬 Chat with Driver</Btn>
      <Btn v="secondary" onClick={()=>go(S.MATCHED)}>←</Btn>
    </div>
    <HI/>
  </div>
);

// ── SCREEN: CHAT ──
const Chat = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"4px 16px 8px",display:"flex",alignItems:"center",gap:10,borderBottom:`1px solid ${C.greyL}`}}>
      <div onClick={()=>go(S.MATCHED)} style={{width:32,height:32,borderRadius:10,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",cursor:"pointer",fontSize:14}}>←</div>
      <div style={{width:36,height:36,borderRadius:12,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:14,color:"#fff",fontWeight:800,fontFamily:F.d}}>JR</div>
      <div style={{flex:1}}><div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy}}>Juan Reyes</div><div style={{fontFamily:F.b,fontSize:10,color:C.green}}>● Online • ETA 5 min</div></div>
      <div style={{fontSize:16,cursor:"pointer"}}>📞</div>
    </div>
    <div style={{flex:1,padding:"12px 16px",overflow:"auto"}}>
      <div style={{textAlign:"center",marginBottom:12}}><span style={{fontFamily:F.b,fontSize:9,color:C.grey,background:C.light,padding:"4px 10px",borderRadius:6}}>Today 2:34 PM</span></div>
      {[
        {from:"driver",msg:"Magandang hapon po! On my way na. I'll be there in about 8 minutes 🚛"},
        {from:"me",msg:"Thank you! I'm at the gas station beside McDonald's EDSA"},
        {from:"driver",msg:"Copy po! I can see the location. White Montero po ba?"},
        {from:"me",msg:"Yes correct! White Montero GLS"},
      ].map((m,i)=>(
        <div key={i} style={{display:"flex",justifyContent:m.from==="me"?"flex-end":"flex-start",marginBottom:8}}>
          <div style={{maxWidth:"75%",padding:"10px 14px",borderRadius:16,...(m.from==="me"?{background:C.g2,borderBottomRightRadius:4,color:"#fff"}:{background:C.white,borderBottomLeftRadius:4,color:C.navy,border:`1px solid ${C.greyL}`})}}>
            <div style={{fontFamily:F.b,fontSize:12,lineHeight:1.4}}>{m.msg}</div>
          </div>
        </div>
      ))}
    </div>
    <div style={{padding:"8px 16px",borderTop:`1px solid ${C.greyL}`,display:"flex",gap:8}}>
      <div style={{display:"flex",gap:6}}>
        {["📍","📷","🎤"].map(e=><div key={e} style={{width:36,height:36,borderRadius:12,background:C.light,display:"flex",alignItems:"center",justifyContent:"center",fontSize:16,cursor:"pointer"}}>{e}</div>)}
      </div>
      <div style={{flex:1,padding:"10px 14px",borderRadius:12,background:C.white,border:`1.5px solid ${C.greyL}`,fontFamily:F.b,fontSize:12,color:C.grey}}>Type a message...</div>
      <div style={{width:36,height:36,borderRadius:12,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:16,cursor:"pointer"}}>➤</div>
    </div>
    <HI/>
  </div>
);
// ── SCREEN: SOS ──
const SOS = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:"linear-gradient(180deg,#1a0000,#4a0000 50%,#8B0000)",alignItems:"center",justifyContent:"center",position:"relative"}}>
    <Bar light/>
    <div style={{position:"absolute",inset:0,background:"radial-gradient(circle at center,rgba(255,71,87,0.15) 0%,transparent 70%)"}}/>
    <div style={{position:"relative",width:140,height:140,marginBottom:24}}>
      <div style={{position:"absolute",inset:0,border:"3px solid rgba(255,71,87,0.2)",borderRadius:70}}/>
      <div style={{position:"absolute",inset:20,background:"linear-gradient(135deg,#FF4757,#FF6B6B)",borderRadius:55,display:"flex",alignItems:"center",justifyContent:"center",fontSize:50,boxShadow:"0 8px 40px rgba(255,71,87,0.5)",cursor:"pointer"}}>🆘</div>
    </div>
    <div style={{fontFamily:F.b,fontSize:10,fontWeight:700,color:"rgba(255,255,255,0.5)",letterSpacing:3,marginBottom:8}}>TAP FOR HELP</div>
    <div style={{fontFamily:F.d,fontSize:22,fontWeight:800,color:"#fff",marginBottom:8}}>Emergency SOS</div>
    <div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.5)",textAlign:"center",padding:"0 40px",lineHeight:1.5}}>This will alert nearby providers, ops center, and PNP-HPG with your live location</div>
    <div style={{display:"flex",flexDirection:"column",gap:8,width:"100%",padding:"24px 24px 0"}}>
      {[{ic:"📞",t:"Call 911",s:"Emergency Services"},{ic:"🚔",t:"Alert PNP-HPG",s:"Highway Patrol Group"},{ic:"📍",t:"Share Live Location",s:"Send to emergency contacts"}].map(a=>(
        <div key={a.t} style={{display:"flex",alignItems:"center",gap:12,padding:"14px 16px",borderRadius:14,background:"rgba(255,255,255,0.08)",cursor:"pointer"}}>
          <span style={{fontSize:22}}>{a.ic}</span>
          <div><div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:"#fff"}}>{a.t}</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)"}}>{a.s}</div></div>
        </div>
      ))}
    </div>
    <div style={{position:"absolute",bottom:40}}><Btn v="ghost" onClick={()=>go(S.MATCHED)} sm st={{color:"rgba(255,255,255,0.4)"}}>Cancel</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: CONDITION REPORT ──
const Condition = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Pre-Tow Condition Report" onBack={()=>go(S.TRACKING)}/>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      <Card elevated st={{background:C.g1,marginBottom:14,border:"none"}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:20}}>📸</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:"#fff"}}>8 photos required before towing</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.6)"}}>All photos are timestamped & GPS-tagged</div></div>
        </div>
      </Card>
      <div style={{display:"grid",gridTemplateColumns:"1fr 1fr 1fr 1fr",gap:6,marginBottom:14}}>
        {["Front","Rear","Left","Right","FL Tire","FR Tire","RL Tire","RR Tire"].map((l,i)=>(
          <div key={l} style={{aspectRatio:"1",borderRadius:12,background:i<3?C.green+"20":C.light,border:i<3?`2px solid ${C.green}`:`2px dashed ${C.greyL}`,display:"flex",flexDirection:"column",alignItems:"center",justifyContent:"center",gap:2}}>
            <span style={{fontSize:i<3?14:18}}>{i<3?"✅":"📷"}</span>
            <span style={{fontFamily:F.b,fontSize:8,fontWeight:600,color:i<3?C.green:C.grey}}>{l}</span>
          </div>
        ))}
      </div>
      <Card st={{marginBottom:12,textAlign:"center",padding:"20px",border:`2px dashed ${C.greyL}`}}>
        <span style={{fontSize:26}}>🎥</span>
        <div style={{fontFamily:F.d,fontSize:12,fontWeight:600,color:C.navy,marginTop:4}}>Record 360° Walk-Around</div>
        <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>Min 30 seconds • GPS & timestamp embedded</div>
      </Card>
      <Card st={{background:C.cream,border:`1.5px solid ${C.gold}33`}}>
        <div style={{display:"flex",alignItems:"center",gap:8}}>
          <span style={{fontSize:14}}>🔒</span>
          <div style={{fontFamily:F.b,fontSize:10,color:"#5D6D7E"}}>All evidence is <b>tamper-proof</b> with SHA-256 hashing and stored for 1 year.</div>
        </div>
      </Card>
    </div>
    <div style={{padding:"12px 16px 4px"}}><Btn full onClick={()=>go(S.COMPLETE)}>Submit Report (3/8 photos)</Btn></div>
    <HI/>
  </div>
);

// ── SCREEN: COMPLETE ──
const Complete = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{flex:1,display:"flex",flexDirection:"column",alignItems:"center",justifyContent:"center",padding:"0 20px"}}>
      <div style={{width:80,height:80,borderRadius:24,background:C.g3,display:"flex",alignItems:"center",justifyContent:"center",fontSize:38,boxShadow:"0 8px 30px rgba(0,137,123,0.3)",marginBottom:16}}>✅</div>
      <div style={{fontFamily:F.d,fontSize:22,fontWeight:800,color:C.navy,marginBottom:4}}>Tow Complete!</div>
      <div style={{fontFamily:F.b,fontSize:12,color:C.grey,marginBottom:24}}>Your vehicle has been delivered safely</div>
      <Card elevated st={{width:"100%",marginBottom:14}}>
        <div style={{display:"flex",justifyContent:"space-between",marginBottom:6}}>
          <span style={{fontFamily:F.b,fontSize:10,color:C.grey}}>Job ID</span>
          <span style={{fontFamily:F.m,fontSize:10,fontWeight:600,color:C.navy}}>TC-2026-00847</span>
        </div>
        {[["Service","Flatbed Tow"],["Distance","12.4 km"],["Duration","45 min"]].map(([l,v],i)=>(
          <div key={i} style={{display:"flex",justifyContent:"space-between",padding:"4px 0"}}>
            <span style={{fontFamily:F.b,fontSize:11,color:C.grey}}>{l}</span>
            <span style={{fontFamily:F.b,fontSize:11,fontWeight:600,color:C.navy}}>{v}</span>
          </div>
        ))}
        <div style={{height:1,background:C.greyL,margin:"8px 0"}}/>
        <div style={{display:"flex",justifyContent:"space-between",alignItems:"center"}}>
          <span style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy}}>Total Paid</span>
          <div style={{textAlign:"right"}}><div style={{fontFamily:F.d,fontSize:18,fontWeight:800,color:C.orange}}>₱1,850</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>via GCash 💚</div></div>
        </div>
      </Card>
      <Card st={{width:"100%",background:C.cream,border:`1.5px solid ${C.gold}33`,marginBottom:20}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:20}}>⭐</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:"#F57F17"}}>+1 Suki Point earned!</div><div style={{fontFamily:F.b,fontSize:10,color:"#795548"}}>1 more booking to reach Gold tier</div></div>
        </div>
      </Card>
      <Btn full onClick={()=>go(S.RATE)}>⭐ Rate Your Experience</Btn>
      <div style={{marginTop:10}}><Btn v="secondary" full onClick={()=>go(S.HOME)}>Back to Home</Btn></div>
    </div>
    <HI/>
  </div>
);

// ── SCREEN: RATE ──
const Rate = ({go}) => {
  const [stars,setStars]=useState(0);
  const [tags,setTags]=useState([]);
  const allTags=["Fast response","Professional","Careful handling","Great communication","Clean truck","Fair price"];
  return (
    <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
      <Bar/><Back title="Rate Your Experience" onBack={()=>go(S.COMPLETE)}/>
      <div style={{flex:1,overflow:"auto",padding:"0 20px"}}>
        <div style={{textAlign:"center",marginBottom:20}}>
          <div style={{width:60,height:60,borderRadius:18,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:26,color:"#fff",fontWeight:800,fontFamily:F.d,margin:"0 auto 10px"}}>JR</div>
          <div style={{fontFamily:F.d,fontSize:16,fontWeight:700,color:C.navy}}>Juan Reyes</div>
          <div style={{fontFamily:F.b,fontSize:11,color:C.grey}}>Flatbed Tow • TC-2026-00847</div>
        </div>
        <div style={{display:"flex",justifyContent:"center",gap:8,marginBottom:20}}>
          {[1,2,3,4,5].map(i=>(
            <div key={i} onClick={()=>setStars(i)} style={{fontSize:36,cursor:"pointer",filter:i<=stars?"none":"grayscale(100%) opacity(0.3)",transition:"all .2s"}}>⭐</div>
          ))}
        </div>
        {stars>0&&(
          <>
            <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8,textAlign:"center"}}>{stars>=4?"WHAT WAS GREAT?":"WHAT COULD IMPROVE?"}</div>
            <div style={{display:"flex",flexWrap:"wrap",gap:6,justifyContent:"center",marginBottom:16}}>
              {allTags.map(t=>(
                <div key={t} onClick={()=>setTags(p=>p.includes(t)?p.filter(x=>x!==t):[...p,t])} style={{padding:"8px 14px",borderRadius:10,border:tags.includes(t)?`2px solid ${C.orange}`:`1.5px solid ${C.greyL}`,background:tags.includes(t)?C.cream:C.white,fontFamily:F.b,fontSize:11,fontWeight:tags.includes(t)?700:500,color:tags.includes(t)?C.orange:C.navy,cursor:"pointer"}}>{t}</div>
              ))}
            </div>
            <div style={{padding:"12px 14px",borderRadius:12,border:`1.5px solid ${C.greyL}`,background:C.white,fontFamily:F.b,fontSize:12,color:C.grey,marginBottom:16}}>Add a comment (optional)...</div>
            <Card st={{background:C.cream,border:`1.5px solid ${C.gold}33`,marginBottom:16}}>
              <div style={{display:"flex",alignItems:"center",gap:8}}>
                <span style={{fontSize:14}}>💡</span>
                <div style={{fontFamily:F.b,fontSize:10,color:"#5D6D7E"}}>Your review helps build a trusted community of Suki providers in your area.</div>
              </div>
            </Card>
          </>
        )}
      </div>
      <div style={{padding:"12px 20px 4px"}}><Btn full onClick={()=>go(S.HOME)} st={{opacity:stars>0?1:0.4}}>Submit Review →</Btn></div>
      <HI/>
    </div>
  );
};

// ── SCREEN: HISTORY ──
const History = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"4px 20px 10px"}}><div style={{fontFamily:F.d,fontSize:18,fontWeight:800,color:C.navy}}>Activity</div></div>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      {[
        {d:"Feb 20",s:"Flatbed Tow",fr:"Makati",to:"Toyota Shaw",p:"₱1,850",st:"complete"},
        {d:"Feb 14",s:"Jumpstart",fr:"BGC",to:"On-site",p:"₱450",st:"complete"},
        {d:"Jan 28",s:"Fuel Delivery",fr:"SLEX Alabang",to:"On-site",p:"₱280",st:"complete"},
        {d:"Jan 15",s:"Tire Change",fr:"Ortigas Ave",to:"On-site",p:"₱650",st:"complete"},
      ].map((h,i)=>(
        <Card key={i} st={{marginBottom:8}}>
          <div style={{display:"flex",justifyContent:"space-between",marginBottom:6}}>
            <span style={{fontFamily:F.b,fontSize:10,color:C.grey}}>{h.d}</span>
            <span style={{fontFamily:F.b,fontSize:8,fontWeight:700,color:C.green,background:`${C.green}15`,padding:"2px 8px",borderRadius:4}}>COMPLETED</span>
          </div>
          <div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy,marginBottom:4}}>{h.s}</div>
          <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>{h.fr} → {h.to}</div>
          <div style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.orange,marginTop:6}}>{h.p}</div>
        </Card>
      ))}
    </div>
    <Nav active="history" go={go}/>
    <HI/>
  </div>
);

// ── SCREEN: PROFILE ──
const Profile = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"8px 20px 16px",display:"flex",alignItems:"center",gap:14}}>
      <div style={{width:56,height:56,borderRadius:18,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:26,color:"#fff",fontWeight:800,fontFamily:F.d}}>D</div>
      <div><div style={{fontFamily:F.d,fontSize:16,fontWeight:800,color:C.navy}}>David</div><div style={{fontFamily:F.b,fontSize:11,color:C.grey}}>+63 9XX XXX XXXX</div><div style={{display:"flex",gap:6,marginTop:4}}><span style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.gold,background:`${C.gold}20`,padding:"2px 8px",borderRadius:4}}>⭐ SUKI SILVER</span></div></div>
    </div>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      {[
        {ic:"🚗",t:"My Vehicles",s:"2026 Montero GLS"},
        {ic:"💳",t:"Payment Methods",s:"GCash • **** 8847"},
        {ic:"📍",t:"Saved Locations",s:"Home, Office, Toyota Shaw"},
        {ic:"🔔",t:"Notifications",s:"Push, SMS, Email"},
        {ic:"🛡️",t:"Safety & Privacy",s:"Emergency contacts, data"},
        {ic:"🌙",t:"Appearance",s:"System default"},
        {ic:"❓",t:"Help & Support",s:"FAQs, chat support"},
        {ic:"📄",t:"Terms & Privacy",s:"Legal documents"},
      ].map((m,i)=>(
        <Card key={i} st={{marginBottom:6,padding:"12px 14px"}}>
          <div style={{display:"flex",alignItems:"center",gap:12}}>
            <span style={{fontSize:18}}>{m.ic}</span>
            <div style={{flex:1}}><div style={{fontFamily:F.d,fontSize:12,fontWeight:600,color:C.navy}}>{m.t}</div><div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>{m.s}</div></div>
            <span style={{fontSize:14,color:C.greyL}}>→</span>
          </div>
        </Card>
      ))}
      <Btn v="ghost" full st={{marginTop:8,color:C.coral}}>Log Out</Btn>
    </div>
    <Nav active="profile" go={go}/>
    <HI/>
  </div>
);

// ── SCREEN: SUKI LOYALTY ──
const Suki = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/>
    <div style={{padding:"4px 20px 10px"}}><div style={{fontFamily:F.d,fontSize:18,fontWeight:800,color:C.navy}}>Suki Rewards</div></div>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      <Card elevated st={{background:C.g1,marginBottom:14,border:"none",padding:20}}>
        <div style={{display:"flex",alignItems:"center",gap:12}}>
          <div style={{width:50,height:50,borderRadius:16,background:"rgba(245,166,35,0.2)",display:"flex",alignItems:"center",justifyContent:"center",fontSize:26}}>⭐</div>
          <div><div style={{fontFamily:F.d,fontSize:16,fontWeight:800,color:"#fff"}}>Silver Member</div><div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.5)"}}>4 of 6 bookings to Gold</div></div>
        </div>
        <div style={{marginTop:14,height:8,borderRadius:4,background:"rgba(255,255,255,0.1)"}}>
          <div style={{width:"66%",height:"100%",borderRadius:4,background:C.g2}}/>
        </div>
        <div style={{display:"flex",justifyContent:"space-between",marginTop:6}}>
          <span style={{fontFamily:F.b,fontSize:9,color:"rgba(255,255,255,0.4)"}}>Silver</span>
          <span style={{fontFamily:F.b,fontSize:9,color:C.gold}}>Gold →</span>
        </div>
      </Card>
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginBottom:8}}>YOUR BENEFITS</div>
      {[["5% off all services","Active now"],["Priority matching","Silver+"],["Free condition report","Silver+"],["10% off + VIP support","Gold tier"]].map(([b,s],i)=>(
        <Card key={i} st={{marginBottom:6,padding:"12px 14px"}}>
          <div style={{display:"flex",alignItems:"center",justifyContent:"space-between"}}>
            <span style={{fontFamily:F.b,fontSize:12,fontWeight:600,color:C.navy}}>{b}</span>
            <span style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:i<3?C.green:C.grey,background:i<3?`${C.green}15`:`${C.grey}15`,padding:"2px 8px",borderRadius:4}}>{s}</span>
          </div>
        </Card>
      ))}
      <div style={{fontFamily:F.b,fontSize:9,fontWeight:700,color:C.grey,letterSpacing:1.5,marginTop:14,marginBottom:8}}>POINTS HISTORY</div>
      {[{d:"Feb 20",a:"+1 point",r:"Flatbed Tow"},{d:"Feb 14",a:"+1 point",r:"Jumpstart"},{d:"Jan 28",a:"+1 point",r:"Fuel Delivery"},{d:"Jan 15",a:"+1 point",r:"Tire Change"}].map((p,i)=>(
        <Card key={i} st={{marginBottom:6,padding:"10px 14px"}}>
          <div style={{display:"flex",justifyContent:"space-between"}}>
            <div><span style={{fontFamily:F.b,fontSize:11,fontWeight:600,color:C.navy}}>{p.r}</span><div style={{fontFamily:F.b,fontSize:9,color:C.grey}}>{p.d}</div></div>
            <span style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.green}}>{p.a}</span>
          </div>
        </Card>
      ))}
    </div>
    <Nav active="suki" go={go}/>
    <HI/>
  </div>
);

// ── SCREEN: TYPHOON MODE ──
const Typhoon = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:"linear-gradient(180deg,#1B3A5C,#0B1D33)"}}>
    <Bar light/>
    <div style={{padding:"8px 20px 0"}}>
      <div style={{display:"flex",alignItems:"center",gap:8,marginBottom:12}}>
        <div style={{padding:"4px 10px",borderRadius:8,background:"rgba(255,71,87,0.2)",fontFamily:F.b,fontSize:10,fontWeight:700,color:C.coral}}>⚠️ TYPHOON ALERT</div>
      </div>
      <div style={{fontFamily:F.d,fontSize:20,fontWeight:800,color:"#fff",marginBottom:4}}>Typhoon Mode Active</div>
      <div style={{fontFamily:F.b,fontSize:11,color:"rgba(255,255,255,0.5)",marginBottom:16}}>Signal #3 — Surge pricing may apply</div>
    </div>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      <Card st={{background:"rgba(255,255,255,0.06)",border:"1px solid rgba(255,255,255,0.1)",marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:22}}>🌊</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:"#fff"}}>Flood Level: Knee-Deep</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)"}}>EDSA-Guadalupe area • Updated 5 min ago</div></div>
        </div>
      </Card>
      <Card st={{background:"rgba(255,255,255,0.06)",border:"1px solid rgba(255,255,255,0.1)",marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:22}}>🚛</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:"#fff"}}>2 Trucks Available</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)"}}>Heavy-duty flatbeds only during typhoon</div></div>
        </div>
      </Card>
      <Card st={{background:"rgba(255,71,87,0.1)",border:"1px solid rgba(255,71,87,0.2)",marginBottom:10}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:22}}>💰</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.coral}}>Surge: 1.5× Base Rate</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)"}}>MMDA-regulated typhoon rate cap applies</div></div>
        </div>
      </Card>
      <Card st={{background:"rgba(0,196,140,0.1)",border:"1px solid rgba(0,196,140,0.2)"}}>
        <div style={{display:"flex",alignItems:"center",gap:10}}>
          <span style={{fontSize:22}}>🛡️</span>
          <div><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.green}}>Safety Guaranteed</div><div style={{fontFamily:F.b,fontSize:10,color:"rgba(255,255,255,0.4)"}}>Full insurance coverage during calamity</div></div>
        </div>
      </Card>
    </div>
    <div style={{padding:"12px 16px 4px",display:"flex",flexDirection:"column",gap:8}}>
      <Btn full onClick={()=>go(S.DIAGNOSE)}>🚛 Book Emergency Tow</Btn>
      <Btn v="ghost" full onClick={()=>go(S.HOME)} st={{color:"rgba(255,255,255,0.4)"}}>← Back to Home</Btn>
    </div>
    <HI/>
  </div>
);

// ── SCREEN: PROVIDERS LIST ──
const Providers = ({go}) => (
  <div style={{flex:1,display:"flex",flexDirection:"column",background:C.bg}}>
    <Bar/><Back title="Nearby Providers" onBack={()=>go(S.HOME)}/>
    <div style={{flex:1,overflow:"auto",padding:"0 16px"}}>
      {[
        {n:"Juan Reyes",r:"4.9",j:"847",d:"0.8 km",eta:"5 min",v:true},
        {n:"Mark Santos",r:"4.7",j:"423",d:"1.2 km",eta:"8 min",v:true},
        {n:"Rico Dela Cruz",r:"4.8",j:"612",d:"2.1 km",eta:"12 min",v:false},
      ].map((p,i)=>(
        <Card key={i} st={{marginBottom:8}}>
          <div style={{display:"flex",alignItems:"center",gap:12}}>
            <div style={{width:46,height:46,borderRadius:14,background:C.g2,display:"flex",alignItems:"center",justifyContent:"center",fontSize:14,color:"#fff",fontWeight:800,fontFamily:F.d}}>{p.n.split(" ").map(w=>w[0]).join("")}</div>
            <div style={{flex:1}}>
              <div style={{display:"flex",alignItems:"center",gap:6}}>
                <span style={{fontFamily:F.d,fontSize:13,fontWeight:700,color:C.navy}}>{p.n}</span>
                {p.v&&<span style={{fontFamily:F.b,fontSize:7,fontWeight:700,color:"#fff",background:C.teal,borderRadius:3,padding:"1px 5px"}}>✓</span>}
              </div>
              <div style={{fontFamily:F.b,fontSize:10,color:C.grey}}>★ {p.r} • {p.j} jobs • {p.d}</div>
            </div>
            <div style={{textAlign:"right"}}><div style={{fontFamily:F.d,fontSize:12,fontWeight:700,color:C.orange}}>{p.eta}</div><div style={{fontFamily:F.b,fontSize:9,color:C.grey}}>ETA</div></div>
          </div>
        </Card>
      ))}
    </div><HI/>
  </div>
);
// ── MAIN APP ──
const SCREENS = {
  [S.LOGO]:LogoShowcase,[S.SPLASH]:Splash,[S.LOGIN]:Login,[S.HOME]:Home,
  [S.DIAGNOSE]:Diagnose,[S.SERVICE]:Service,[S.VEHICLE]:Vehicle,[S.DROPOFF]:Dropoff,
  [S.PRICE]:Price,[S.MATCHING]:Matching,[S.MATCHED]:Matched,[S.TRACKING]:Tracking,
  [S.CHAT]:Chat,[S.CONDITION]:Condition,[S.COMPLETE]:Complete,[S.RATE]:Rate,
  [S.SOS]:SOS,[S.PROVIDERS]:Providers,[S.HISTORY]:History,[S.PROFILE]:Profile,
  [S.TYPHOON]:Typhoon,[S.SUKI]:Suki,
};

export default function TowCommandApp() {
  const [screen,setScreen]=useState(S.LOGO);
  const Comp=SCREENS[screen]||Home;
  return (
    <div style={{width:375,height:812,margin:"0 auto",borderRadius:44,overflow:"hidden",boxShadow:"0 20px 60px rgba(11,29,51,0.3),0 0 0 1px rgba(11,29,51,0.1)",display:"flex",flexDirection:"column",fontFamily:F.d,position:"relative",background:C.bg}}>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@400;500;600;700;800&display=swap');
        @keyframes logoPulse{0%,100%{transform:scale(1)}50%{transform:scale(1.05)}}
        @keyframes dotPulse{0%,100%{opacity:.3;transform:scale(.8)}50%{opacity:1;transform:scale(1.2)}}
        @keyframes spin{0%{transform:rotate(0)}100%{transform:rotate(360deg)}}
        * { box-sizing:border-box; -webkit-font-smoothing:antialiased; }
        ::-webkit-scrollbar{display:none}
      `}</style>
      <Comp go={setScreen}/>
    </div>
  );
}
