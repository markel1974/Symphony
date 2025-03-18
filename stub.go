package main

const treeStub = `{
 "c64": {
  "children": {
   "cartridges_c64": {
    "children": {
     "easyFlash:0": {
      "details": {
       "id": "easyFlash:0"
      }
     }
    },
    "details": {
     "id": "cartridges_c64"
    }
   },
   "dynamic_throttle": {
    "details": {
     "id": "dynamic_throttle"
    }
   },
   "iec": {
    "children": {
     "c1541:8": {
      "children": {
       "mos6510": {
        "details": {
         "id": "mos6510"
        }
       },
       "mos6522:1": {
        "details": {
         "id": "mos6522:1"
        }
       },
       "mos6522:2": {
        "details": {
         "id": "mos6522:2"
        }
       },
       "pic_6510": {
        "details": {
         "id": "pic_6510"
        }
       },
       "pla_c1541": {
        "details": {
         "id": "pla_c1541"
        }
       },
       "roms_c1541": {
        "details": {
         "id": "roms_c1541"
        }
       }
      },
      "details": {
       "id": "c1541:8"
      }
     }
    },
    "details": {
     "id": "iec"
    }
   },
   "joystick_c64:1": {
    "details": {
     "id": "joystick_c64:1"
    }
   },
   "joystick_c64:2": {
    "details": {
     "id": "joystick_c64:2"
    }
   },
   "keyboard_c64": {
    "details": {
     "id": "keyboard_c64"
    }
   },
   "mos6510": {
    "details": {
     "id": "mos6510"
    }
   },
   "mos6526:1": {
    "children": {
     "timer:A": {
      "details": {
       "id": "timer:A"
      },
      "properties": {
       "cnt": false,
       "countMode": 0,
       "cr": 0,
       "crNew": 0,
       "crNewPending": false,
       "timer": 65535,
       "timerLatch": 65535,
       "timerLatchLow": 0,
       "timerState": 0,
       "toggleMode": false
      }
     },
     "timer:B": {
      "details": {
       "id": "timer:B"
      },
      "properties": {
       "cnt": false,
       "countMode": 0,
       "cr": 0,
       "crNew": 0,
       "crNewPending": false,
       "timer": 65535,
       "timerLatch": 65535,
       "timerLatchLow": 0,
       "timerState": 0,
       "toggleMode": false
      }
     },
     "tod": {
      "details": {
       "id": "tod"
      }
     }
    },
    "details": {
     "id": "mos6526:1"
    }
   },
   "mos6526:2": {
    "children": {
     "timer:A": {
      "details": {
       "id": "timer:A"
      },
      "properties": {
       "cnt": false,
       "countMode": 0,
       "cr": 0,
       "crNew": 0,
       "crNewPending": false,
       "timer": 65535,
       "timerLatch": 65535,
       "timerLatchLow": 0,
       "timerState": 0,
       "toggleMode": false
      }
     },
     "timer:B": {
      "details": {
       "id": "timer:B"
      },
      "properties": {
       "cnt": false,
       "countMode": 0,
       "cr": 0,
       "crNew": 0,
       "crNewPending": false,
       "timer": 65535,
       "timerLatch": 65535,
       "timerLatchLow": 0,
       "timerState": 0,
       "toggleMode": false
      }
     },
     "tod": {
      "details": {
       "id": "tod"
      }
     }
    },
    "details": {
     "id": "mos6526:2"
    }
   },
   "mos6569": {
    "details": {
     "id": "mos6569"
    }
   },
   "mos6581": {
    "details": {
     "id": "mos6581"
    },
    "properties": {
     "ad1": 0,
     "ad2": 0,
     "ad3": 0,
     "cr1": 0,
     "cr2": 0,
     "cr3": 0,
     "env3": 0,
     "fcHI": 0,
     "fcLO": 0,
     "freqHI1": 0,
     "freqHI2": 0,
     "freqHI3": 0,
     "freqLO1": 0,
     "freqLO2": 0,
     "freqLO3": 0,
     "modeVol": 0,
     "osc3": 0,
     "potX": 255,
     "potY": 255,
     "pwHI1": 0,
     "pwHI2": 0,
     "pwHI3": 0,
     "pwLO1": 0,
     "pwLO2": 0,
     "pwLO3": 0,
     "resFilt": 0,
     "sr1": 0,
     "sr2": 0,
     "sr3": 0
    }
   },
   "pic_6510": {
    "details": {
     "id": "pic_6510"
    }
   },
   "pla_c64": {
    "children": {
     "ports": {
      "details": {
       "id": "ports"
      }
     }
    },
    "details": {
     "id": "pla_c64"
    }
   },
   "quartz": {
    "details": {
     "id": "quartz"
    }
   },
   "roms_c64": {
    "details": {
     "id": "roms_c64"
    }
   }
  },
  "details": {
   "id": "c64"
  }
 }
}
`
