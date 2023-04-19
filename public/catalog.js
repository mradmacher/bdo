const processes = {
  R1: "Wykorzystanie głównie jako paliwa lub innego środka wytwarzania energii",
  R2: "Odzysk/regeneracja rozpuszczalników",
  R3: "Recykling lub regeneracja substancji organicznych, które nie są stosowane jako rozpuszczalniki (w tym kompostowanie i inne biologiczne procesy przekształcania)",
  R4: "Recykling lub odzysk metali i związków metali",
  R5: "Recykling lub odzysk innych materiałów nieorganicznych",
  R6: "Regeneracja kwasów lub zasad",
  R7: "Odzysk składników stosowanych do redukcji zanieczyszczeń",
  R8: "Odzysk składników z katalizatorów",
  R9: "Powtórna rafinacja oleju lub inne sposoby ponownego użycia olejów",
  R10: "Obróbka na powierzchni ziemi przynosząca korzyści dla rolnictwa lub poprawę stanu środowiska",
  R11: "Wykorzystywanie odpadów uzyskanych w wyniku któregokolwiek z procesów wymienionych w pozycji R 1 – R 10",
  R12: "Wymiana odpadów w celu poddania ich któremukolwiek z procesów wymienionych w pozycji R 1 – R 11",
  R13: "Magazynowanie odpadów poprzedzające którykolwiek z procesów wymienionych w pozycji R1 – R 12 (z wyjątkiem wstępnego magazynowania u wytwórcy odpadów)",
  D1: "Składowanie w gruncie lub na powierzchni ziemi (np. składowiska itp.)",
  D2: "Przetwarzanie w glebie i ziemi (np. biodegradacja odpadów płynnych lub szlamów w glebie i ziemi itd.)",
  D3: "Głębokie zatłaczanie (np. zatłaczanie odpadów w postaci umożliwiającej pompowanie do odwiertów, wysadów solnych lub naturalnie powstających komór itd.)",
  D4: "Retencja powierzchniowa (np. umieszczanie odpadów ciekłych i szlamów w dołach, poletkach poletkach osadowych lub lagunach itd.)",
  D5: "Składowanie na składowiskach w sposób celowo zaprojektowany (np. umieszczanie w uszczelnionych oddzielnych komorach, przykrytych i izolowanych od siebie wzajemnie i od środowiska itd.)",
  D6: "Odprowadzanie do wód z wyjątkiem mórz i oceanów",
  D7: "Odprowadzanie do mórz i oceanów, w tym lokowanie na dnie mórz",
  D8: "Obróbka biologiczna, niewymieniona w innej pozycji niniejszego załącznika, w wyniku której powstają ostateczne związki lub mieszanki, które są unieszkodliwiane za pomocą któregokolwiek spośród procesów wymienionych w poz. D 1 – D 12",
  D9: "Obróbka fizyczno-chemiczna, niewymieniona w innej pozycji niniejszego załącznika, w wyniku której powstają ostateczne związki lub mieszaniny unieszkodliwiane za pomocą któregokolwiek spośród procesów wymienionych w pozycjach D 1 – D 12 (np. odparowanie, suszenie, kalcynacja itp.)",
  D10: "Przekształcanie termiczne na lądzie",
  D11: "Przekształcanie termiczne na morzu",
  D12: "Trwałe składowanie (np. umieszczanie pojemników w kopalniach itd.)",
  D13: "Sporządzanie mieszanki lub mieszanie przed poddaniem odpadów któremukolwiek z procesów wymienionych w pozycjach D 1 – D 12",
}

const codes = {
  '01': {
    desc: 'Odpady powstające przy poszukiwaniu, wydobywaniu, fizycznej i chemicznej przeróbce rud oraz innych kopalin',
    entries: {
      '01': {
        desc: 'Odpady z wydobywania kopalin',
        entries: {
          '01': {
            desc: 'Odpady z wydobywania rud metali (z wyłączeniem 01 01 80)',
          },
          '02': {
            desc: 'Odpady z wydobywania kopalin innych niż rudy metali',
          },
          '80': {
            desc: 'Odpady skalne z górnictwa miedzi, cynku i ołowiu',
          },
        },
      },
      '03': {
        desc: 'Odpady z fizycznej i chemicznej przeróbki rud metali',
        entries: {
          '04*': {
            desc: 'Odpady z przeróbki rud siarczkowych powodujące samoczynne zakwaszenie środowiska w czasie składowania',
          },
          '05*': {
            desc: 'Inne odpady poprzeróbcze zawierające substancje niebezpieczne (z wyłączeniem 01 03 80)',
          },
          '06': {
            desc: 'Inne odpady poprzeróbcze niż wymienione w 01 03 04, 01 03 05, 01 03 80 i 01 03 81',
          },
        },
      },
      '04': {
        desc: 'Odpady z fizycznej i chemicznej przeróbki kopalin innych niż rudy metali',
        entries: {
          '07': {
            desc: 'Odpady zawierające niebezpieczne substancje z fizycznej i chemicznej przeróbki kopalin innych niż rudy metali',
          },
          '08': {
            desc: 'Odpady żwiru lub skruszone skały inne niż wymienione w 01 04 07',
          },
          '09': {
            desc: 'Odpadowe piaski i iły',
          },
        }
      },
    }
  },
  '02': {
    desc: 'Odpady z rolnictwa, sadownictwa, upraw hydroponicznych, rybołówstwa, leśnictwa, łowiectwa oraz przetwórstwa żywności',
    entries: {
      '01': {
        desc: 'Odpady z rolnictwa, sadownictwa, upraw hydroponicznych, leśnictwa, łowiectwa i rybołówstwa',
        entries: {
          '01': {
            desc: 'Osady z mycia i czyszczenia',
          },
          '02': {
            desc: 'Odpadowa tkanka zwierzęca',
          },
          '03': {
            desc: 'Odpadowa masa roślinna',
          },
        }
      },
    },
  },
  '03': {
    desc: 'Odpady z przetwórstwa drewna oraz z produkcji płyt i mebli, masy celulozowej, papieru i tektury',
    entries: {
      '01': {
        desc: 'Odpady z przetwórstwa drewna oraz z produkcji płyt i mebli',
        entries: {
          '01': {
            desc: 'Odpady kory i korka',
          }
        }
      },
    }
  }
}
