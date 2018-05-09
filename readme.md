metaGraf
========

**metaGraf** provides a generic and implementation agnostic
"structure" of metadata about a software component. **metagraf**
is inspired by the <a href="https://12factor.net">twelve-factor app</a> guidelines to 
aid automation tasks or decisions about a component or collection of compoenents.

**metaGraf** operates on an individual or collections of metagraph(s)
 (software components) to produce aggregated metadata to support your
toolchain or pipelines with information that can be acted upon.

One of the goals of **metaGraf** is to indentify missing nodes or edges
(components) when comparing a running enviroment with the graph/branch
of a component not currently deployed or a new version of an existing
component. Desired state vs existing state.

In other words, determining what needs to be present or changed to 
fulfil the explicit dependencies of the new component entering an
environment.

Another goal is to aid in documentation of software components and
their dependencies on a component level.

Background
-
metaGraf is currently a research project and a place to experiment
with a structure for describing software components and how that
information can be used to assist CI and CD pipelines, developers,
architects, operations and a organization as a whole.

I have not found many projects that solve the complexities of
managing software components in an enviroment similar to the goals
of metaGraf.

The <a href="https://getnelson.github.io/nelson/">Nelson</a> project is
the closest thing I have found (after some pointers). It has a concept of
topology map (graph) of deployments in an environment.

<a href="http://ddd.ward.wiki.org/view/welcome-visitors/view/ward-cunningham">Ward Cunningham</a> 
is dabbling with something in this space at a broader scope: 
http://ddd.ward.wiki.org/view/about-the-el-dorado-project/

If anyone is  interested in this subject, please reach out and hopefully
 we can get a discussion going. Input and suggestions are always welcome.
  

Direction
-
Since cloud-native now eats the world, the goal is to enable building 
Kubernetes Operators/Controllers that act on the metadata and 
collections of metadata. The structure so far is also inspired by a 
Kubernetes resource so a metaGraf could be a CRD. 

mgraf
-
A little tool to help communicate what metagraf attempt to solve and basis for
further discussion.

Usage:

> mgraf -c /path/to/collection/of/metagraphs 

Usage straight from source:

> go run mgraf -c /path/to/collection/of/metagraphs

You can use the example collection provided to experiment. It produces output like ihis if the resulting

<img src="data:image/png:base64,iVBORw0KGgoAAAANSUhEUgAAAbQAAAG7CAYAAACrYmCTAAAABmJLR0QA/wD/AP+gvaeTAAAgAElE
                                QVR4nOzdd1RU19rH8S+gSLWCiF2siA1UQKwI9t6wxMSCNWo0iUaN3Whs0VhTLOg1igZNVLDFHlFs
                                KFjATrBiwUpvM+8fllcFpQ2cAZ7PWq5175R9fofMmWfOPvvsraNWq9UIIYQQOZyu0gGEEEIITZCC
                                JoQQIleQgiaEECJXyKd0ACFyGrVazYMHDwgLCyMiIgK5DC3SS19fHxMTEypUqICpqanScXINKWhC
                                pIFKpeLkyZMcPHiQM6dPExEZqXQkkUuULlUSp4aNaNOmDeXLl1c6To6mI6Mchfg0Pz8/flmxnPth
                                D6hVrSz1a1fAumJJSloUxsTIAF1dHaUjihwmITGJlxExhN59zPkrdzhx7ib3Hz7FyakBX345glKl
                                SikdMUeSgibER9y7d4/Fi3/m7NlzNHGoRt/OTpQsXljpWCIXUqvVnL0YisdWX+4/fEaPHm7069cP
                                fX19paPlKFLQhEjBuXPnmD5tGuZFjRn2mTM2leUXs8h6SSoVuw+fZ8M2P8pVsGLWrNkULiw/otJK
                                CpoQH9i5cydLliymUb0qjB7QEv38cqlZZK+7YU+ZsXQHKvIxZ+48ubaWRlLQhHjHwYMHmT17Nr07
                                OtKnYwN0dOT6mFBGRGQsPyzfwYPwaH759VeKFy+udCStJwVNiNeuXr3K6NFf0bZZLQb1bKp0HCGI
                                iY1n7Jw/yadvwtJlyzE0NFQ6klaTG6uFAF6+fMnECROoY12WgT2aKB1HCAAMDfSZOqoTjx8/ZOHC
                                n5SOo/WkoAkBeHh4oEMS44a0kWH4QqtYmBXkG/dWHDx4iMDAQKXjaDUpaCLPCw0NZefOnQzo3ghD
                                AxkmLbRPvZoVsK9dkSWLfyYpKUnpOFpLCprI81atWknFcsVxbmCtdBQhPmpQr6bcvXePQ4cOKR1F
                                a0lBE3na48ePOXnyFF1b1ZURjUKrlbIoQgPbSvh4eysdRWtJQRN52vHjxzEokJ8GdpWVjiJEqpwb
                                WHMpKIhnz54pHUUrSUETeVpAwDlqVStDPj05FIT2q1O9HHp6ujI45CPkKBZ5WsjNm1QsKzesipyh
                                gH4+SpcoRkhIiNJRtJIUNJGnhT95gnmxdK5HlfSCExtX07v7MKxqdaNQ5W5YOg6j6cDFTF/vx6m7
                                MaiyJm7OE7oFB6sOmHbdxp1s22gCoUe86NW8O8aV5+OTwVYigg8wbuBoKtfuRtEan1O37zJ+PfuS
                                9M9EoZk8bxQrYsKTJ08y2UruJAVN5GlxcfEU0M+fjnc8Z9s3Y2g57QTqFu5471vHg0sbCPT6mhHV
                                n7Bhxhyad/LgeJYlTqfoQL52dqOC+75sLCj/L3DLQS4BqsD9bLie9duLuX2aWYNH4Dzfn5DwuAy3
                                ExW4kZbdl7HbpDWb93tyz/cHvit7hUm9xvK1b0S253mXYYF8xMTEaKSt3EYKmsjT1Gp1+kY3XtjO
                                RJ9wzNxGsn5ofWqUMMWwgCEWZaviNnYyHp9ZZF3YjFCrUKnVqNWqDJxZZJLqMn9si6S2TTHgDhu2
                                Xs3iDEl4L1xNcN0R+Pm408wkg82oQlk4/k8uFHRlxfw21C1eAMMi5en5w7eMKBfGqgkbOBafjXk+
                                IKNxP04KmhDpEHntNneA8lalSH4LdgGatKuPefbH+jhjO5Yc2UKoR2vKZvOm43z386eeM0t/bE4F
                                IHTbAXyz9J5gPTrOW47nsNpY6mW8laTTe/G4rsayjTPNDN5t3oo+HcpB2CFWH0rL2ZZm8oi0k4Im
                                RDqYmBXGGAg+FsCDlF7gMJTQsyNonM25tE8sO7cco0R3V+rVdKVvNSD8KOsPp+nUJsM0MdPL5RMX
                                eAzY1qyU7LlqtSphRCyHj1/Ntjwi7WShJyHSw74hnc33s/Hob7gMuM/4oa3pam+JSSo/DdVPr7J6
                                mRfrDlzm6sMYdEzNqG7vyJBRvfisuvGrF+1bSMFhR3h1ElOf3w604daizXj6hXL3eTwpndzU/OZX
                                To4sDUlH6Vl5ATvfPNF6HFGdz7zX3prLU+lV4J1Mz67zvxVbWHsgiOCwGPIVK46VVRVad25J//Y1
                                KPPO2Uma8r/r+XH+OFiaL8aXA+Cz7tbMnnWZHVuO87OrM6kOw3m5j3Z1lnEkrfv6i6YmlFZz9eZ9
                                wIRSlgbJntWxKEYJIOS/e7ygFoU0tFWhGXKGJkR6GNVl/oovcLGE0H//ZnifIZR1GEn7r9ew1CeY
                                Oyldq3/kx8BO4xmzK5ZOP8znSqAnFzcPxvX5YYZ0m8CMc6+7r1p+y8uQ7axpAXCTOZMOUbzvWI4d
                                38TNv3thp1eDJWc38ktTfdCtyOwjPq++4AH0mvBniAdz7IzptGTTqy/499r7wONTDO78HaN8Ymg7
                                9UcunNvMde/vmer4gjXfTWSAZ3j687/jvvcBjtq50qvMq/9fpksLmuaD6EMH2JqWAXoFW7IrJB37
                                qjExvHiZBBhglNJKLcaGGAO8jERubdY+UtCESKfC9XrgfXgV+3/6nEEtqlIi9g6Hd2xn4ujx2DT6
                                hjE+9/j/jrV49ixYgdc9fTpOHs93zUpjZmRIicr2TFn6Bc3Uofw03YfkdxUl0HToGIY4WlLUUB/z
                                Op/he30Og4oUpNdgFyxUN/ll9QUS3nlH4tkdrHjQnNFtUxt9EM/eBcvYdCc/nadOYHzzclga62Ni
                                VpZWI8cxo2m+916b/vwP8Nx6DdfuTXh7h1+RRnzhYgBJF1m/7VEa/9Ka2FcNU6tfD2zRQYZmaB8p
                                aEJkhL4ZTl3dWPL7TwQHbiRg4xgmdiiP8bPrrPpmNvOC3nQQBuOz7yXoWtOuecH32zCvjXMVUF06
                                ye5kF+TKU692gQ8fBKCAU2eG2ehw769tbH17mhDF1pUHKD2wEw6pDkAIxvufF0A1Wjb7sLvQhH5r
                                t3FgoFnG8185wIb/HPi89bvFxpCO3Z0ohJrTWw5w7d12rnlSx6oDxu/8qzDtoob2Nb0MKVRQD4gl
                                OqWz7ehYogEKGlNY05sWmSbX0ITILD0TqjRwYXKDprQt9RWNf7vDX3tuMcXGCuKf8ygC4BxDa3Vg
                                aIoN3Ofmf0CJdx8zwNjoYxssyeBBDvz09SmWbLhD71Fl4L89LDlty8RFabht4E2mAoUxT+HyV4qv
                                TXN+NX5bDnE96jG9bHxTbvP6Qf53vjeza78+x6nSh8CQPh8JkMl9TTcdqlYsCdzhXlgs8P51NPXD
                                JzwAzCqUkutnWkgKmhDpcXY1FYcl8L8zw2mU7Ml82DlaY/LbHZ69iHz1kH4RLAoCUU5suDKRLho6
                                oyjSviufzz/Jb+u3cWjol7Dah6hek2mfWoEC0C9McVMg4jmPo4BPvSe9+RPP88d2NaO2ejPX7sNO
                                OTUnZg3C1eMRm7ZeYEbt2mn6AsrUvmaAtWNNzJfdIfDSTehq895zVy/eJJoCtGtYNWs2LjJFuhyF
                                SA+1GtWTAHYFpDT8XM31iyFEokttm3KvH7OmU6tCkHSNY/7J33Pjt28wbbiSk+m9P0vPmlEDqqH3
                                5AiLV29jsU9JRvWrnMYDujodWxUCrvDP4agPngtjfvuOWM+6/Hr6rvTljzq4n7+LuvBFsmIGoEMD
                                NxeqAA+997MvNjv2Nf30HNowsLIO9/cc4ei7412SQtm0MxQsXXB3Trk7WChLCpoQ6RbG8pE/MNP7
                                MjceRhEbH0v43ZvsXLWAHktvYFyjJ9O7vumQ0qfld1/Rt+xzVo9fwNIjtwmLiCfmeRjHPBfjtuwZ
                                Xb/vimMGztzK9+5CR9MEDi78g8CWXehbIvX3vMnUetwoepeJZ8esucw/fJsHUfFEPLjOn1PmseBR
                                faYPtH795ZCe/FFs33KKaj1cqf6xTVdxpV9tHYg4wfp/PiymWbGvn/DiDCObu2HRZD5bHr7zuG55
                                vp3rRq3n+/jyu72cexxH7LNbeE1byLL/LBk0py+NC6ShHZHtdNRqdbbPiCOEtnB2dmbC8A40rl8l
                                bW9QxRIa4I/PvlPsPxPKf2HhhIXHoDIwpbRVRZq0bsWY/g2o9MGQb/XzG6xf4YXH/ksE349Bt2BR
                                KtjY8dkQN4Y3NH/V9RbgQeVu27j/3jud2BAykS4phlHjP3coTVfqMnH3r0yu9sFZ0Xv3tb1SoNMk
                                nv7s+Ordz66zboUXa/cHERwWS76i5tg4NOXrr7vTvvz7NwSnmv/BDpo5rebMmzfUGci1v7tQ6t1G
                                7m6jURMPAt59zKwd+08Pw+mjf/A07usbhxZTdNBBUprHw3XOZnb0fKef8vkpvuz8E1sT67H8r/G4
                                fXBJ7mXQfmYu9GG7/12eqowoX9ueQd/0Y3jdQu+PcPxUO+nJk0Zzf92JnrEl06ZNS/d7czspaCJP
                                S3dBE0JhUtA+TrochRBC5ApS0IQQQuQKUtBEnqanp0eSSpbjFDlHkkqFrq58dadE/ioiTzM2MiI6
                                WjMLLwqRHaKi4zExyeYpv3IIKWgiT7O0LMH9hzLNrMg57j54hqWlpdIxtJIUNJGnVa5SlSshcvOQ
                                yBnCn0Xw5NlLKleurHQUrSQFTeRp9vb2XA25x/OX0UpHESJVpwJDMDAoQK1atZSOopWkoIk8zd7e
                                HmMjY/b5XlQ6ihCp+ufoJZo0aUr+/PmVjqKVpKCJPK1AgQK079CB7fsDiIxK6+SCQmQ/v7PXCbn9
                                kC5dUp43RkhBE4K+ffuSL18BPL1PKB1FiBQlJCax7q/jtGjhSrVq1ZSOo7WkoIk8z8jIiIHu7uw8
                                FMjNW2ldTVmI7LPZ5yRPX0QxZEjKK9KJV6SgCQG0bt2aOrXrMHPZDp4+T/ss8EJktWP+1/hz5ymG
                                DRtOsWLFlI6j1aSgCQHo6uoyfcYMjIwLMXPZDmLjEpSOJATX/nvAotV76dKlCx07dlQ6jtaTgibE
                                ayYmJsyZO5fHT6OZMH+LnKkJRflf+I9JP22lTh1bRowYoXScHEGWjxHiA/fv32fihAlER71g6qhO
                                VCxXXOlIIg9Rq9X4HAxk1ebDtGzZim+//ZZ8+fIpHStHkIImRAoiIyOZMX06AYEBtHOuw2edGmBi
                                bKB0LJHLhdx+xO+eRwi6fpdBgwbRp08fpSPlKFLQhPgIlUrF3r17WbVyJWp1Ip1b2NKiUQ2KFEr/
                                KsNCfMr10AfsPHSeQ35BWFtb89VXo6lSRRadTS8paEKkIioqio0bN7LTx4fIqCiqWpXEumIJLC0K
                                Y2JkgK6ujtIRRQ6TkJDEi8gYbt0N58LVuzx8/JwK5cvRu89nuLq6oqMjn6mMkIImRBrFxcVx+vRp
                                zpw5w7WrVwgLe0BkVBSqXLieWlJSEiEhIZQrVw59ff1s335kZGSuXiJFXz8/piYmlC9fgeo2Njg5
                                OckN0xogBU0Ikcz48eP57bffCAoKonTp0tm67cOHD+Pq6sqZM2ews7PL1m2LnE0KmhDiPefPn6d+
                                /fosXbqUYcOGKZKhYcOG6Ovrc/jwYUW2L3ImKWhCiLeSkpJwcHBAX1+fY8eOoaurzK2qvr6+NGnS
                                hL1799KqVStFMoicRwqaEOKtBQsWMGXKFM6dO0f16tUVzdKxY0du3bpFQECAYoVV5CzyKRFCABAa
                                GsqMGTOYNGmS4sUMYP78+QQHB/PHH38oHUXkEHKGJoQAoFWrVty+fZvAwEAKFCigdBwAhgwZwu7d
                                u7l27RpGRkZKxxFaTs7QhBD873//48CBA6xevVprihnAjBkzeP78OcuXL1c6isgB5AxNiDwuPDyc
                                6tWr06tXL5YuXap0nGSmTJnC8uXLuXHjhiyfIj5JCpoQeVyfPn04fvw4ly5dwtTUVOk4yURGRlKp
                                UiX69u3LTz/9pHQcocWky1GIPGzPnj1s2rSJZcuWaWUxg1fL+kyePJlly5YREhKidByhxeQMTYg8
                                Kjo6mho1auDg4MCmTZuUjvNJCQkJVK9enQYNGrB+/Xql4wgtJWdoQuRR33//PS9evGDx4sVKR0lV
                                /vz5mTlzJp6enly5ckXpOEJLyRmaEHnQ6dOncXJyYtWqVQwYMEDpOGmiUqmoU6cONWrUwNPTU+k4
                                QgtJQRMij0lMTKR+/foULlyYQ4cO5ailSry8vOjduzcBAQHUqlVL6ThCy0hBEyKPmTVrFj/++CMX
                                LlygUqVKSsdJF7Vaja2tLVWqVMHLy0vpOELLyDU0IfKQa9euMXv2bKZPn57jihmAjo4OU6dOZevW
                                rZw/f17pOELLyBmaEHmEWq3G1dWV8PBw/P39yZ8/v9KRMkStVuPg4EDp0qX5+++/lY4jtIicoQmR
                                R6xatYp///0XDw+PHFvM4NVZ2uTJk9m+fTtnzpxROo7QInKGJkQe8ODBA6pXr467uzsLFixQOo5G
                                ODg4YGFhgbe3t9JRhJaQgiZEHtCtWzfOnTvHxYsXMTExUTqORuzevZt27dpx8uRJHBwclI4jtIAU
                                NCFyuZ07d9KhQwf++ecfWrZsqXQcjWrcuDGmpqbs3r1b6ShCC0hBEyIXe/nyJTY2Nri4uLBu3Tql
                                42jcvn37aNWqFb6+vjRq1EjpOEJhUtCEyMWGDx/O1q1bCQ4OxtzcXOk4WaJx48YYGxuzd+9epaMI
                                hckoRyFyqRMnTrBy5UqWLFmSa4sZwOTJk/nnn384e/as0lGEwuQMTYhcKC4uDjs7O8qVK5cnri/V
                                q1cPKysrmT0kj5MzNCFyodmzZxMaGsqKFSuUjpItxo8fz19//UVwcLDSUYSC5AxNiFzmypUr1KlT
                                h3nz5jF69Gil42QLlUpFzZo1cXR0ZM2aNUrHEQqRgiZELqJSqWjSpAnx8fGcOHECPT09pSNlGw8P
                                D4YNG8aNGzcoW7as0nGEAqTLUYhcZMWKFZw6dYrff/89TxUzgL59+1KiRAl+/vlnpaMIhUhBEyKX
                                uH37NpMmTWL8+PHY2toqHSfb6evrM3r0aNasWcPz58+VjiMUIF2OQuQSnTp14sqVK5w/fx4DAwOl
                                4ygiIiKCMmXKMGnSJMaNG6d0HJHN5AxNiFxg8+bN+Pj48Ouvv+bZYgZgamrKwIEDWbJkCQkJCUrH
                                EdlMCpoQOcjo0aNZsmQJKpXq7WNPnz5lzJgxDB48mObNmyuYTjuMHj2ahw8fsmXLFqWjiGwmXY5C
                                5CBlypTh7t27NGrUCA8PDypXrszAgQPZs2cPwcHBFClSROmIWqFnz57cuHFDZg/JY6SgCZFD3L9/
                                n1KlSgG8XaDziy++wMPDgy1bttCtWzcl42mVU6dO4ejoyJEjR2jatKnScUQ2kS5HIXIIX19fdHR0
                                AEhISCAhIYF169ZRvHhxSpcurXA67eLg4ICjoyPLly9XOorIRlLQhMghjh079vbM7I2kpCSePHlC
                                gwYNGDJkCBEREQql0z4jRoxg+/bt3Lt3T+koIptIQRMihzh8+DDx8fHJHk9MTEStVrN27VqqVavG
                                nj17FEinfdzc3ChWrBirVq1SOorIJlLQhMgBXrx4weXLl1N93dOnT+Wm4tf09fUZOHAgK1eulCH8
                                eYQUNCFyAD8/v/eG6n8of/78lCxZktOnT9O7d+9sTKbdhg0bxqNHj9i2bZvSUUQ2kIImRA5w/Phx
                                9PX1U3xOV1cXFxcXAgMDqVmzZjYn025ly5alffv2/PLLL0pHEdlACpoQOcDhw4eTdZvp6uqio6PD
                                uHHj2LVrl9yD9hHDhw/n6NGjXL16VekoIovJfWhCaLn4+HgKFixIXFzc28fy5cuHoaEhXl5etG7d
                                WsF02k+lUlGhQgX69OnDnDlzlI4jspCcoQmh5fz9/ZMVs6pVqxIYGCjFLA10dXXp168f69atIzEx
                                Uek4IgtJQRNCy/n6+r69/0xHR4cePXpw+vRprKysFE6WcwwYMIBHjx7xzz//KB1FZCEpaEJouaNH
                                j5KQkEC+fPlYvnw5np6eGBkZKR0rR6lQoQJNmzZl7dq1SkcRWUiuoYlMSUhIIDQ0lGfPnhEdHa10
                                nFxHpVLRqVMn8ufPzw8//ICNjY3SkTTOyMiIIkWKUL58+WQzoWjSH3/8waBBg7h79y7m5uZZth2h
                                HCloIt0iIiLYt28fvr5HuXQpiKSkJKUj5VpRUVFcu3YNGxubjw7bzy309PSoUcOGxo2b0LJlS0xN
                                TTXafnR0NJaWlsyaNYtRo0ZptG2hHaSgiTSLi4tj06ZNbN68CT1dHRxtK1GvZnkqli1OsSImGBrk
                                7i9cJYQ9fIq5WSHy6ekpHSXLxMTG8+RZJDdvP8L/YignA26QpFLTq1dvevfuTYECBTS2rS+++IKQ
                                kBCOHTumsTaF9pCCJtLE19eX5cuWERn5kl7tHWjrXFsKmMgSMbHx7D58ns07T2FiYsrIUV/RuHFj
                                jbS9a9cuOnToQEhICOXLl9dIm0J7SEETn6RWq1mzZg2enp64NrKhf7fGFC4oAxJE1nv+Mpp1f/ly
                                4FgQffr0wd3d/e3yORmVkJCApaUl48ePZ9y4cRpKKrSFFDTxUXFxcfz442z8/Pz4ql8LXBrmvgEJ
                                QvsdPB7E0nX7cWroxPffT8p0F+TgwYMJCAjA399fQwmFtpBh+yJFKpWKH3+czbmz/swe212KmVCM
                                S0MbZo/rzrmz/vz44+xPTtKcFj179uTs2bNcv35dQwmFtpCCJlLk4eGBn58fk75sT40qshqyUFaN
                                KqWZ+lUnTvidwMPDI1NtOTs7Y2FhwZ9//qmhdEJbSEETyfj6+uLp6clX/VpQy7qs0nGEAMCmcilG
                                9XPF09MTX1/fDLejp6dH9+7dpaDlQlLQxHvi4uJYvmwpro1spJtRaB2Xhja4NrJh+bKl781vmV49
                                e/bk0qVLBAUFaTCdUJoUNPEeT09PIiMj+KJLQ6WjCJGiAd2bEBUVyaZNmzLcRqNGjShbtqycpeUy
                                UtDEWxEREfz552Z6tXegaGETpeMIkaJCpob0bGfP5s2biIiIyFAbOjo6dO/eHS8vLw2nE0qSgibe
                                2rdvH3q6OrR1rq10FCE+qa1zbfR0ddi/f3+G2+jSpQtXr17lypUrGkwmlCQFTbzl63sUR9tKMgOI
                                0HqGBvo42lbi6NF/M9yGk5MTFhYW7NixQ4PJhJKkoAng1arIQUHB1KtZXukoQqRJ3RrlCAoKJiEh
                                IUPv19XVpW3btlLQchEpaAKAW7dukZiYiFUZWVZD5AwVyxYnMTGR27dvZ7iNTp06cerUKR48eKDB
                                ZEIp+ZQOILTD06dPATArVjB9b0x6wYnNW1i6zZ9T1x7zJAaMiplTpXo1nJvZ06a5LfVLG+bBX05n
                                GVJjOhtTXCJOh3wGxpQoXwnXbj2Y0r8WJbJsMv0EQo9sY8JML3zu2LP5+nd0SMe7k17c4eDOQ2za
                                7Y9f8H0exuanWAlzKtd2ZNDQjnS1Nk3nf9vM5XmXWdFXy8s8efKEihUrZqiNli1bYmBgwM6dOxk0
                                aFAGkwhtkfe+Z0SKYmJiADDQT89vnOds+2YMLaedQN3CHe9963hwaQOBXl8zovoTNsyYQ/NOHhzP
                                msjpFx3I185uVHDfx50s31hdVl7yIWqnG1YALb7hRYgPUSHePL+4nrPr3HEmmHWzp9BqziVisyBB
                                zO3TzBo8Auf5/oSEZ+yerUNzxtFl+nFUru547/uDB+dWcWB+B0pf/pt+HUczbN+LbM3zLoMCrxYD
                                zczCsoaGhri6ukq3Yy4hBU28J12zmV/YzkSfcMzcRrJ+aH1qlDDFsIAhFmWr4jZ2Mh6fWWRd0IxQ
                                q1Cp1ajVKpSbkVuH/MaFqWTvyopZbSmDihub9nAgUdPbScJ74WqC647Az8edZpm4C8Pc7Us8BtSh
                                qrkRBoamVKjfkt8Xd6Fy0mM2zt3BpWzO80ZmZ95/o2PHjhw8eJCoqCiNtCeUIwVNZFjktdvcAcpb
                                lSL5uMgCNGlXH626Imdsx5IjWwj1aI02TOilZ1WKCgAxT3mQsdupPtU6Hectx3NYbSwz0Z3ZYu5m
                                QmfX4cMm9KyrUtcAuB3GzTT9OtBMnqzQoUMH4uPj2bdvn9JRRCZJQRMZZmJWGGMg+FgAKV5SdxhK
                                6NkRaGZpxtwnKeQeIQBFy1KtiObbz9LbL6IjeBoLVClH9TSeKGnr7SDFixfHwcEBb29vpaOITJJB
                                ISLj7BvS2Xw/G4/+hsuA+4wf2pqu9paYpPIzSf30KquXebHuwGWuPoxBx9SM6vaODBnVi8+qG796
                                0b6FFBx2hCQA6vPbgTbcWrQZT79Q7j6Pf/34+2p+8ysnR5aGpKP0rLyAnW+eaD2OqM5n3mtvzeWp
                                9HpnWS31s+v8b8UW1h4IIjgshnzFimNlVYXWnVvSv30NyhikM//H956E6JfcCjrDwmm7uWtYin6z
                                3Wj0qbe83Ee7Oss4ktZ9/aVJKhky78nuYxzFhK4j2lI5y7eW9Tp16sSCBQtISkpCT0/LTiFFmskZ
                                msg4o7rMX/EFLpYQ+u/fDO8zhLIOI2n/9RqW+gRzJyaF9zzyY2Cn8YzZFUunH+ZzJdCTi5sH4/r8
                                MEO6TWDGudeDBVp+y8uQ7axpAXCTOZMOUbzvWI4d38TNv3thp1eDJWc38ktTfdCtyOwjPq++4AH0
                                mvBniAdz7IzptGTTqy/499r7wONTDO78HaN8Ymg79UcunNvMde/vmer4gjXfTWSAZ3j6839o/yIK
                                WXXA2KojhWv0pXbPJfwdb8vERVNY2KrYp//OBVuyKyQd+5rVwv0YPf8slj2/ZWm7dI6K1VIdOnQg
                                PDycU6dOKR1FZIIUNJEphev1wPvwKvb/9DmDWlSlROwdDu/YzsTR47Fp9A1jfO4R//bV8exZsAKv
                                e/p0nDye75qVxszIkBKV7Zmy9AuaqUP5abrPq2649yTQdOgYhjhaUtRQH/M6n+F7fQ6DihSk12AX
                                LFQ3+WX1Bd69vTbx7A5WPGjO6LapjT6IZ++CZWy6k5/OUycwvnk5LI31MTErS6uR45jRNN97r81Y
                                ft4Z5ehD1PUt/PfvIn7tEM+6L7/CYfherqQ6KEQT+6oBzy4wof9izjp+zUglqWwAACAASURBVO5Z
                                9ciCnlJFWFtbU65cObmOlsNJQROZp2+GU1c3lvz+E8GBGwnYOIaJHcpj/Ow6q76ZzbygNx2Ewfjs
                                ewm61rRr/sEve/PaOFcB1aWT7E52Qa489WoX+PBBAAo4dWaYjQ73/trG1mdvHo1i68oDlB7YCYdU
                                e4+C8f7nBVCNls0+7C40od/abRwYaJbJ/B/QM6B4mcp0Hf09P3c04OY/vzN83f1Xz13zpI5VB4zf
                                +Vdh2kUN7WsmRV9jRr9ZeFcaxf6fm1E2l/XMtWzZkn/++UfpGCIT5Bqa0Cw9E6o0cGFyg6a0LfUV
                                jX+7w197bjHFxgrin/MoAuAcQ2t1YGiKDdzn5n9AiXcfM8DY6GMbLMngQQ789PUplmy4Q+9RZeC/
                                PSw5bcvERWm4beBNpgKFMU/t8leG83+MAY3srWBHIGf8gogcVBKTKn0IDOnzkddncl8zIymMNSNn
                                sr7EYPYvbEzpXFbMAFq1aoWHhwdPnjyhWLFUuoGFVpIzNJFxZ1dTsf6vHEvxyXzYOVpjAjx7Efnq
                                If0iWBQE9JzYcP1191uyf5783CB9MYq078rnlmourt/GofhEDq32IapXV9qnVqAA9AtT3BSIe87j
                                1G5DyoL86tdD3tUxcaTlVuNM7WuGvWD399P4Ia4HPstbYPW2mIUwpVl/JgZk5bazT4sWLdDV1eXA
                                gQNKRxEZJAVNZJxajepJALsC4lN6kusXQ4hEl9o25V4/Zk2nVoUg6RrH/JO/58Zv32DacCUnUxrC
                                +Cl61owaUA29J0dYvHobi31KMqpf5TR+uKvTsVUh4Ar/HP6wooUxv31HrGddRpUl+WM5fubVFbdy
                                tSqTpnOCTO1rRsThv2QmI643469VnaiunSPvNaJgwYLY29tLt2MOJgVNZFIYy0f+wEzvy9x4GEVs
                                fCzhd2+yc9UCeiy9gXGNnkzvWuj1a/Vp+d1X9C37nNXjF7D0yG3CIuKJeR7GMc/FuC17Rtfvu+KY
                                ge6s8r270NE0gYML/yCwZRf6pqnL71Wm1uNG0btMPDtmzWX+4ds8iIon4sF1/pwyjwWP6jN9oPXr
                                A0VD+ZNieXz3OtuWzOFr75folmjEbPeq2bCvn/DiDCObu2HRZD5bHr55UE3o1nl0W3KNR4GbaFLj
                                /Wt7xlajWfThvMAptpNztGrVir1796JWKzeXjMg4uYYmMs7ucw5vscZn3yn2r/uFLXPCCQuPQWVg
                                SmmrijQZ8z1j+jeg0ru/6ovZ89v2BTit8MJj+gR+uB+DbsGiVLCxY8DK+Qxv+HoARoAHlbtt49VQ
                                iTO4W3fAHSc2hEykS0pZjBswprcl21bqMmRQfQw/fP69+9petfdlp0k8/dkRzB1YtX0BDVd4sXb6
                                BOaHxZKvqDk2Dk1Z69Wd9iUzkP/DyYn3L6KQ1aJX/1tHjwJGJliWLUfjISMZPqgF9um5ZJPavr5x
                                aDFFBx18pyvTl15WvgC4ztnMjp7v9FOqVahUatQq9TvTgiVwbM8Z3rlpIXUptpOBPApp3bo1U6dO
                                5dKlS9SsWVPpOCKddNTyU0QAR44cYcaMGezy+FbpKEKkWbuBC5k2bRrNmjXTSHsqlQpLS0vGjRvH
                                2LFjNdKmyD7S5SiEEK/p6uri4uIi19FyKCloQgjxjtatW+Pr6yuz7+dAUtAEwNv561Qq6YEWOUOS
                                6tXYU03PvdiiRQvi4uLw8/PTaLsi60lBEwAYG7+6IB8dk/mFF4XIDtHRr26dMDHR7JRflpaWVKlS
                                hcOHD2u0XZH1pKAJ4NVBDHD3wbNUXimEdrj78Cnw/59dTXJ2dpaClgNJQRMAlChRAlMTEy7fvK90
                                FCHS5MrNMExNTLCw0Py0X87Ozvj7+/Py5UuNty2yjhQ0Abxazr6+vT1nzv+ndBQh0uR0YAj2Dg7o
                                6KRxhdF0cHZ2JikpiWPHUp7YTWgnKWjiLRcXFy5cuc39R8+VjiLEJ917+IyLV+/QvHnzLGm/ePHi
                                VK9eXbodcxgpaOItR0dHSlqWYMN2Gd0ltNvGHScoaVkCR0fHLNuGXEfLeaSgibd0dXX5csRIjp66
                                wqVrd5WOI0SKLt+8z9FTV/hyxEh0dbPuK8zZ2ZnAwECePZOBUjmFFDTxHicnJ+rVrcvvnkeIT0h1
                                GWUhslV8QiK//HGIenXr4uTklKXbatasGWq1mqNHj2bpdoTmSEETyYweM4bHT6NYsnafzDoutIZa
                                rWbJ2n08fhrF6DFjsnx7RYsWpWbNmtLtmINIQRPJlCpViukzZnDM/xqe3ieUjiMEAJ7eJzjmf43p
                                M2ZQqlSpbNlm8+bNOXLkSLZsS2SeFDSRIjs7O0aPHsMm75Os3nxEpsQSilGp1KzefIRN3icZPXoM
                                dnZ22bbtxo0bc/HiRZ4/l5G/OYEsHyM+6eDBg8yfPw/b6uUYN6QNhga5eMlioXViYuNZsHIPAcG3
                                +O678bi4uGTr9h8/foyFhQW7du2iTZs22bptkX5S0ESqgoKCmDJ5Mjok0b9bQ5o7Vc+Sm1mFeEOt
                                VnPIL5h1fx1HjR4/zJqFjY2NIlmqVq1Kjx49mDVrliLbF2knBU2kSUREBGvWrMHHx4dK5S3o2qoe
                                DewqkU9Peq2F5iQmqThx7gZ//+PPjdCHdOjQAXd3d0xNTRXL5O7uzs2bN+VaWg4gBU2kS0hICB4e
                                azhx4iQGBfJTq1oZKpYtjllRUwwN8isdT+RAMbEJhD+N4ObtR1y4cofYuAQaNHBk4EB3rKyslI6H
                                h4cHI0eO5Pnz5+jrS5e7NpOCJjLk8ePH+Pn5ce7cOUJu3uDps2dER8coHSsZlUpFSEgIhQoVwtzc
                                XJEMbya4LViwoCLb13ZGRoYUKVyYipUqY2dnR8OGDTEzM1M61ltXr16lWrVqnDhxIktnJhGZJwVN
                                5Fo3b96ke/fu3Lp1i/Xr19O+fXtFcri5uQHg5eWlyPZF5qjVakqUKMG4ceMYO3as0nHEJ8gFEJEr
                                +fj4UL9+fXR0dDhz5oxixUzkfDo6OjRs2JDjx48rHUWkQgqayFUSExOZPn06nTt3pn379hw7doyK
                                FSsqHUvkcA0bNuTYsWMyc46Wy6d0ACE05dGjR/Tp04fjx4+zcuVK3N3dlY4kcomGDRsSHh7OtWvX
                                qFq1qtJxxEfIGZrIFf7991/q1KnDnTt3OH36tBQzoVF2dnYYGBjg5ydLK2kzKWgiR1Or1SxZsgRX
                                V1ccHBw4ffo0NWvWVDqWyGX09fWxs7Pj1KlTSkcRnyAFTeRYT548oX379owdO5ZZs2axbds2ChUq
                                pHQskUs5ODhIQdNyUtBEjnTu3Dnq16/PhQsXOHr0KOPHj1c6ksjlHBwcuHjxIlFRUUpHER8hBU3k
                                OCtXrsTJyYkKFSrg7+9PgwYNlI4k8gAHBweSkpI4d+6c0lHER0hBEzlGZGQkffr0YdiwYYwZM4b9
                                +/djYWGhdCyRR5QvX54SJUpIt6MWk2H7Ike4evUq3bt358GDB+zZs4dWrVopHUnkQfb29lLQtJic
                                oQmtt3HjRurVq4ehoSH+/v5SzIRiZGCIdpOCJrRWXFwco0ePpm/fvvTp04djx45Rrlw5pWOJPMzB
                                wYE7d+5w7949paOIFEhBE1rpzp07NGvWjLVr1/Lnn3/y+++/y9IdQnH29vbo6elx+vRppaOIFEhB
                                E1pn165d1KlThxcvXnDy5Mm3s9ULoTRTU1OqVq3KmTNnlI4iUiAFTWiNpKQkpk+fTseOHWnXrh3+
                                /v5Ur15d6VhCvKdu3boydF9LSUETWuHx48e0adOGefPmsWjRItavX4+RkZHSsYRIxtbWlrNnzyod
                                Q6RAhu0Lxfn6+tKrVy+MjIw4efIktWvXVjqSEB9lZ2dHeHg4d+/epXTp0krHEe+QMzShmDcTC7u4
                                uFCvXj3OnDkjxUxoPVtbW3R0dAgICFA6iviAFDShiJcvX+Lm5sbYsWP54Ycf2L59O4ULF1Y6lhCp
                                KliwIFZWVnIdTQtJl6PIdgEBAfTo0YPY2FiOHDlCw4YNlY4kRLrY2trKGZoWkjM0ka3Wr19Po0aN
                                KFOmDP7+/lLMRI4kBU07SUET2SI2NpbBgwfTv39/Ro0axYEDByhRooTSsYTIEDs7O27fvk14eLjS
                                UcQ7pMtRZLlr167Ro0cP7t27x65du2jTpo3SkYTIFDs7OwACAwNxdXVVOI14Q87QRJbavn079vb2
                                5M+fnzNnzkgxE7lC8eLFKVmypAwM0TJS0ESWSExMZMKECXTp0oWePXvi5+dHhQoVlI4lhMbY2dnJ
                                dTQtI12OQuPu3r1Lz549uXDhAps2baJXr15KRxJC42xtbfHy8lI6hniHnKEJjTp06BD16tXj6dOn
                                nDx5UoqZyLVsbW25fv06ERERSkcRr0lBExqhVquZN28eLVq0wNXVFX9/f2xsbJSOJUSWsbOzQ6VS
                                cf78eaWjiNekoIlMCw8Pp02bNkybNo1FixaxYcMGjI2NlY4lRJYqV64cZmZmch1Ni8g1NJEpZ86c
                                wc3NDZVKxb///ouDg4PSkYTINrVr15aCpkXkDE1k2MqVK2nUqBE1atQgICBAipnIc+zs7GTovhaR
                                MzSRbhEREQwaNIi//vqLyZMnM3XqVHR15beRyHtsbW1ZvHgxYWFhXL9+naCgIC5dusSDBw/466+/
                                lI6X50hBE+ly+fJlunfvzuPHj9m7d6/MkiDylJcvXxIcHMylS5cICgrixIkTAJQsWRKA/Pnzk5SU
                                JL0VCpGCJt76999/SUxMxMXFJcXn169fz/Dhw6lXrx4HDhzA0tIymxMKoRy1Wk2zZs0ICAhAT08P
                                PT094uPj33tNQkIC+fPnp06dOgqlzNukn0gAEB0dTb9+/XBzc+PevXvvPRcbG8vo0aPp378/gwYN
                                kmIm8iQdHR0WLFgAQFJSUrJi9i65ZUUZUtAEAFOmTOHevXu8fPmSrl27kpCQAMCtW7do0qQJ69at
                                Y8uWLSxZsoT8+fMrnFYIZbi4uNC6detPHgMJCQnUqFEjG1OJN6SgCU6dOsXixYtJTEwkMTGRc+fO
                                MWHCBLy9vbG1tSUxMZGAgAC6deumdFQhFPfzzz+TlJT0yddUr149m9KId0lBy+Pi4uL44osv0NHR
                                eftYYmIiixYtonPnzri5uXHixAmsrKwUTCmE9qhWrRru7u4fPUsrUqQI5ubm2ZxKgBS0PG/mzJnc
                                vHkz2S9OHR0d9PX1GTNmDAUKFFAonRDaadasWR8taNLdqBwpaHnY+fPnmTdvXordJ2q1mqSkJLp0
                                6UJ0dLQC6YTQXsWLF2f8+PHky/f+QHF9fX0Z4aggKWh5VGJiIp9//vl7XY0pvebGjRuMHDkyG5MJ
                                kTOMGzcOMzOz944hlUolIxwVJAUtj5o7dy7BwcEkJiZ+8nWJiYmsXbuWtWvXZlMyIXIGQ0NDZs+e
                                /d5jiYmJUtAUJAUtDwoKCmLmzJkfHamlo6PztivF2tqaWbNm0bBhw+yMKESO0L9/f2xsbNDT03v7
                                mIxwVI7MFJLHqFQqBgwYkOzxN0UsISGBypUr07t3b3r37k3VqlUVSClEzqCrq8vChQtp1aoVAGZm
                                ZhQtWlThVHmXFLQ8ZtGiRfj7+6NWq9HR0UFXVxeVSoWdnR29e/emW7dulC9fXumYQuQYLVu2xNXV
                                lQMHDlCzZk2l4+RpGitoCQkJhIaG8uzZMxkVp6Xu3r3L999//7aYVatWDRcXF5o0afL2vpnQ0FBC
                                Q0MVy6ivr4+JiQkVKlTA1NRUsRzi4yIiIvjvv/+IjIz85PRPeUnv3r05dOgQxYoV48iRI0rHyRGM
                                jIwoUqQI5cuX19jsQzpqtVqd0TdHRESwb98+fH2PculSUKp3zwvlqNXqt0vFm5ubY25ujr6+vsKp
                                Pq10qZI4NWxEmzZtcvRZo5ubGwBeXl4KJ8m40NBQdu/ezQm/49y9d1/pOFrp6tWrFCxYUOY5TSc9
                                PT1q1LChceMmtGzZMlM/ZDNU0OLi4ti0aRObN29CT1cHR9tK1KtZnopli1OsiAmGBtr9RZkXRURG
                                ExuXgHmxQkpH+aSExCReRsQQevcx56/c4cS5m9x/+BQnpwZ8+eUISpUqpXTEdMvJBe3evXv88ssK
                                /PxOUNKiKA3sKlK7WhnKlzanoKkh+fPppd5IHhH28CnhT19Q07qC0lFyhJjYeJ48i+Tm7Uf4Xwzl
                                ZMANklRqevV6df0+IxM6pLug+fr6snzZMiIjX9KrvQNtnWtLARNZRq1Wc/ZiKB5bfbn/8Bk9erjR
                                r18/rT+7fFdOLGjx8fH873//Y8sWL0paFGVg90bUrVn+k/ctCpEZMbHx7D58ns07T2FiYsrIUV/R
                                uHHjdLWR5oKmVqtZs2YNnp6euDayoX+3xhQuaJSh4EKkV5JKxe7D59mwzY9yFayYNWs2hQsXVjpW
                                muS0gvb8+XMmT57Erf9C6NvFibbOtdGTFclFNnn+Mpp1f/ly4FgQffr0wd3dPc0/pNJU0OLi4vjx
                                x9n4+fnxVb8WuDSUGweFMu6GPWXG0h2oyMecufNyxLW1nFTQQkNDmThhArokMO2rTpS2lCHoQhkH
                                jwexdN1+nBo68f33k9LUBZnqzy6VSsWPP87m3Fl/Zo/tLsVMKKq0ZVEWTepNsUL6jP32Wx49eqR0
                                pFzj0aNHjP32W4oVys+iSb2lmAlFuTS0Yfa47pw768+PP85GpVKl+p5UC5qHhwd+fn5M+rI9NaqU
                                1khQITLD1MSAGWO6YGqcj0nfTyQmJkbpSDleXFwc06ZOxdhAl2lfdcbUxEDpSEJQo0pppn7ViRN+
                                J/Dw8Ej19Z8saL6+vnh6evJVvxbUsi6rsZBCZJahgT5TR3Xi8eOHLFz4k9JxcrwFC+YTFnaX6WO6
                                YGwkywUJ7WFTuRSj+rni6emJr6/vJ1/70YIWFxfH8mVLcW1kI92MQitZmBXkG/dWHDx4iMDAQKXj
                                5FiBgYEcPHiIb9xbYWFWUOk4QiTj0tAG10Y2LF+2lLi4uI++7qMFzdPTk8jICL7oIpPSCu1Vr2YF
                                7GtXZMnin+XG/gxQqVQsW7oUhzqVqFdT7p8S2mtA9yZERUWyadOmj74mxYIWERHBn39upld7B4oW
                                NsmygEJowqBeTbl77x6HDh1SOkqOc/DgQW7fuY17zyZKRxHikwqZGtKznT2bN28iIiIixdekWND2
                                7duHnq4ObZ1rZ2lAITShlEURGthWwsfbW+koOY73jh00sK1EKYsiSkcRIlWv7onUYf/+/Sk+n2JB
                                8/U9iqNtJZkBROQYzg2suRQUxLNnz5SOkmM8ffqUoOBgnBtYKx1FiDQxNNDH0bYSR4/+m+LzyQpa
                                fHw8QUHB1KtZPquzCaExdaqXQ09PVwaHpENgYCB6errUqV5O6ShCpFndGuUICgomISEh2XPJCtqt
                                W7dITEzEqox5toQTQhMK6OejdIlihISEKB0lxwgJCaF0iWIU0JdlEUXOUbFscRITE7l9+3ay55J9
                                kp8+fQqAWbF0Dt9NesGJzVtYus2fU9ce8yQGjIqZU6V6NZyb2dOmuS31Sxumfid3rqYi/MJhVvxx
                                lL2nbvLf4ygSDQpiWbwYpa0q0qiRLS6Na1O/rDGanMM86cUdDu48xKbd/vgF3+dhbH6KlTCncm1H
                                Bg3tSFdr03T+d0kg9Mg2Jsz0wueOPZuvf0eHDCXTVDuvFCtiwpMnTzLRQt7y5MkTzIpmYNCXHOsf
                                cZYhNaaz8cPlIHXyYWRqjHnJstjVrUPHnq3oWqNQlqyurKljXdPfGRHBB5j5kw/bz97lSZIRFerY
                                M+jrfgyrW5D0TndtVvTV8jJPnjyhYsWK7z2XLNObWRcM0vWr7TnbvhlDy2knULdwx3vfOh5c2kCg
                                19eMqP6EDTPm0LyTB8fTGTzLRAfytbMbFdz3cSe7tql6xv4546jZbT2nijVm/trl3DjvRej++Wya
                                0oo6kWeZN2UOLs1+5m8Njz4/NGccXaYfR+Xqjve+P3hwbhUH5neg9OW/6ddxNMP2vUhzWzG3TzNr
                                8Aic5/sTEv7x+0Gyq513GRbIJ7OGpENsbGw6j3OQY/1T6rLykg9RO92wAmjxDS9CfIi69ic3D8xn
                                3dcNsAjZw5COg2jw3UGuZsFHVVPHuia/M6ICN9Ky+zJ2m7Rm835P7vn+wHdlrzCp11i+9k15tOKn
                                GBR4tRhoSgtJf7TIpmuZiAvbmegTjpnbSNYPrU+NEqYYFjDEomxV3MZOxuMzi3SHzlJqFSq1GrVa
                                RYZXN02XeM4smo7bqns4/jAPnwmuNK5YGBP9/JiYWVCjUSvmrJvL4mZZt0KzuduXeAyoQ1VzIwwM
                                TalQvyW/L+5C5aTHbJy7g0tpaiUJ74WrCa47Aj8fd5pl+I4OTbXzPlnaJP3S/TeTYz399PQpaFYS
                                +xYdWLhxKTuHliFk62Jajz7I7SwIpZljXUPtqEJZOP5PLhR0ZcX8NtQtXgDDIuXp+cO3jCgXxqoJ
                                GziWzkXPP/WZ1UivQOS129wByluVIvm4yAI0aVcfrboiZ2zHkiNbCPVoTXZM6KW+uZ1Rv4WQUKMb
                                C3uVSLk7Ua84/Ue6kBWzZbaYu5nQ2XWSbVfPuip1DYDbYdxM04GlR8d5y/EcVhvLTPWJaqodkd3k
                                WM8sUxp99x1T6ujx6MBKJuxM/xnKp2jqWNdUO0mn9+JxXY1lG2eavTs9qJ4VfTqUg7BDrD6kmR4a
                                0FBBMzErjDEQfCyABym9wGEooWdHkL6l2nKPU5v2cFEFtm0bveqK+AhdO3euhkymR3Z9yUdH8DQW
                                qFKO6mn8oa6pWznklpCcSY51DdApgfsX9chPNN5/HOZhdmwzA8e6Jtq5fOICjwHbmpWSPVetViWM
                                iOXw8auZCPQ+zVyXtG9IZ/P9bDz6Gy4D7jN+aGu62ltikkq5VD+9yuplXqw7cJmrD2PQMTWjur0j
                                Q0b14rPqxq9etG8hBYcd4dVlpfr8dqANtxZtxtMvlLvP40npclPNb37l5MjSkHSUnpUXsPPNE63H
                                EdX5zHvtrbk8lV7vzMWqfnad/63YwtoDQQSHxZCvWHGsrKrQunNL+revQZl3fmWkKT+POH46HDCg
                                RjXLdP5hU/ByH+3qLONIWvf3l4/PAPFk9zGOYkLXEW2pnPlkIi+QY/0Tx3ramdazoQanCAgIwi+x
                                I11S+ibWwmM9fe2ouXrzPmBCKcvkqzfoWBSjBBDy3z2gViZS/T/NDEQyqsv8FV/gYgmh//7N8D5D
                                KOswkvZfr2GpTzB3Urr4+ciPgZ3GM2ZXLJ1+mM+VQE8ubh6M6/PDDOk2gRnnXp+GtvyWlyHbWdMC
                                4CZzJh2ieN+xHDu+iZt/98JOrwZLzm7kl6b6oFuR2Ud8Xv0HB9Brwp8hHsyxM6bTkk2v/oO/194H
                                Hp9icOfvGOUTQ9upP3Lh3Gaue3/PVMcXrPluIgM8w9OfnyeEPQIwoagmJmMo2JJdIenY348J92P0
                                /LNY9vyWpe1kQlqRRnKsf+JYTwfzwlgAJD3jwdOPvEbbjvV0txPDi5dJgAFGhik8bWyIMcDLyIxn
                                +oDGRtYWrtcD78Or2P/T5wxqUZUSsXc4vGM7E0ePx6bRN4zxucf/X/uLZ8+CFXjd06fj5PF816w0
                                ZkaGlKhsz5SlX9BMHcpP031IfkdRAk2HjmGIoyVFDfUxr/MZvtfnMKhIQXoNdsFCdZNfVl/g3dvt
                                Es/uYMWD5oxum9rIg3j2LljGpjv56Tx1AuObl8PSWB8Ts7K0GjmOGU3zvffa9OfXQXNjFjK5v88u
                                MKH/Ys46fs3uWfWQSY9Eesixnpb8qVCTxkEqWnKsZ8V3hlr9+m+gucFcmr1VRN8Mp65uLPn9J4ID
                                NxKwcQwTO5TH+Nl1Vn0zm3lBbzoNgvHZ9xJ0rWnX/INKb14b5yqgunSS3ck66ctTr3bKazUVcOrM
                                MBsd7v21ja1vZz+KYuvKA5Qe2AmHVK9LBeP9zwugGi2bfdiFYEK/tds4MNAsA/mLYVkcIIInH/sl
                                9jHXPKlj1QHjd/5VmHYxc/sbfY0Z/WbhXWkU+39uRlkZlCEyQo71VPKn4tHTV9cg8xXFsijafaxn
                                uB1DChXUA2KJTunMPTqWaICC6e+y/ZismyJAz4QqDVyY3KApbUt9RePf7vDXnltMsbGC+Oc8igA4
                                x9BaHRiaYgP3ufkfUOLdxwwwNvrYBksyeJADP319iiUb7tB7VBn4bw9LTtsycVEahhK/yVSgMOap
                                /X3Tlb84jRzM4dJjLl0Jg2bpuI5WpQ+BIX0+8mQG9jcpjDUjZ7K+xGD2L2xMaSlmQhPkWE8h/6e9
                                9A8mCNCxs6FhPrT3WM9UOzpUrVgSuMO9sFjg/eto6odPeACYVSiVwXDJaeYM7exqKtb/lWMpPpkP
                                O0drTIBnL173leoXwaIgoOfEhus+RIWk9M+TnxukL0aR9l353FLNxfXbOBSfyKHVPkT16kr7tPwA
                                0C9McVMg7jmPo1J7bfry2/dph60eBOw5TuhHG1XjP3cIxhVHseBmVuzvC3Z/P40f4nrgs7wFVm8/
                                mCFMadafiQFp26bI4+RYz3x+VRir/zhLIsZ0+bxZmm5zUOZYz3w71o41MQcCLyX/Urt68SbRFMC5
                                YdW0BkqVZgqaWo3qSQC7AlK6Q07N9YshRKJLbZs3k6Ba06lVIUi6xjH/5O+58ds3mDZcycn0zpih
                                Z82oAdXQe3KExau3sdinJKP6VU7jTlanY6tCwBX+OfzhpzyM+e07Yj3rMqoM5Nep0JHlo6tgcPEv
                                xm559LqN98WH7GSyZxhF2/VkSMUUXpCp/Y3Df8lMRlxvxl+rOlFdRsyLjJJjPZP5Izi2YAGzzydR
                                otVQ5rRN42QK2X6sa6YdPYc2DKysw/09Rzj67tiZpFA27QwFSxfcnVPuWs4IDV5DC2P5yB+Y6X2Z
                                Gw+jiI2PJfzuTXauWkCPpTcwrtGT6V0LvX6tPi2/+4q+ZZ+zevwCPfpSqwAAIABJREFUlh65TVhE
                                PDHPwzjmuRi3Zc/o+n1XHDNwmly+dxc6miZwcOEfBLbsQt80dwPo03rcKHqXiWfHrLnMP3ybB1Hx
                                RDy4zp9T5rHgUX2mD7R+/QdLb/781Bk5lS3DynBi0jg6zj/Isf9eEJUQz4sHdznqtYaOfVZxodrn
                                eM9tRKGPZszI/qoJ3TqPbkuu8ShwE01qvN9Pb2w1mkUfzvH54gwjm7th0WQ+WzJzk4ym2hFaRo71
                                dOVXJRD5JIwz+3cy9rOv6PD7HSr2+Jq9PztTOh3jIbLvWNfgd4Zueb6d60at5/v48ru9nHscR+yz
                                W3hNW8iy/ywZNKcvjTVXzzR0Dc3ucw5vscZn3yn2r/uFLXPCCQuPQWVgSmmrijQZ8z1j+jeg0rtV
                                vpg9v21fgNMKLzymT+CH+zHoFixKBRs7Bqycz/CGry/KBnhQuds27gNwBnfrDrjjxIaQiXRJKYtx
                                A8b0tmTbSl2GDKpPstGi793r8qq9LztN4unPjmDuwKrtC2i4wou10ycwPyyWfEXNsXFoylqv7rQv
                                mYH8bxXC+bu5nG99hOX/O8Q3fdfw3+NoEguYUrZKNVoMm8nqz+pQOn86//ap7S8JHNtzhvAU3vpR
                                ahUqlRq1Sp18JNahxRQddJD//7HlSy8rXwBc52xmR09jzbcjtIcc65841j+YnHj/IgpZLQIdPQxN
                                TDAvWQa7em1YObEVXTIyOXG2Heua/c4wtu3LP1stmLnQBzeXlTxVGVG+tj2zNvdjeF3NTveno1ar
                                39v+kSNHmDFjBrs8vtXohoTIanN/3YmesSXTpk1TOsp73NzcAPDy8lI4yftmzJhBUlQYE4a3VzqK
                                EOnSbuBCpk2bRrNmzd57PO+t8CCEECJXkoImhBAiV0hW0PT0Xl3dVKkUW2xBiAxJUqnQ1ZXfaGml
                                q6tLkiqlMbdCaK83n9k3tepdyY5+Y+NXF+OjYzQ3pb8Q2SEqOh4TEw0trpYHGBsbEx2TzsWohFBY
                                dPSrz2xKx3qygmZp+Womi7sPniV7sRDa7O6DZ28/vyJ1lpaW/9fencdFVe9/HH/NwgwwbIqA4I4b
                                phbuimnuWhkubaap3TK1rDS71yXzppY/U8uttFxvi1mpaW5lrrihKIgL7ooLKrKo7MvAzPz+QEkU
                                ZOewfJ6PBw9gzjnf8x5H5jPnnO/5fuXvXJQ51yMyxhDM7m/9kYJWtWpV7O3sOHPpZvEnE6KIRN+N
                                5/bdOOrXl4lw8qpBgwZE34kj+m7RTjIpRHE6eykcezs73NweHfbrkYKmUqlo1bo1R45fLpFwQhSF
                                gGOhWFvrefLJoplXqSJo2rQp1tZ6Ao7lcaw1IUqBw8dCad2mDapspi/J9gp6165dOXH2GjcjY4o9
                                nBBF4e+9IXTs+AxWVvm9M73i0ul0dOz4DH/vPaV0FCHy5EbEXU6eC6NLly7ZLs+2oLVt2xYP96qs
                                /MO/WMMJURT8gy4Qei2Cfv2yHU9CPEa/fv0IvRaBf9AFpaMIkaufNxzEw70qbdu2zXZ5tgVNrVbz
                                7qj32BtwlpDz14s1oBCFkZZu4vvfD9C9eze8vLyUjlPmeHl50a1bN5av3osxLV3pOELk6Mylm+wN
                                OMu7o97L8facHG/a8fHxoWWLFixe5Sf/0UWp9eumQ9yJTWT48OxnqhK5GzFiBLEJKfy2OUDpKEJk
                                y5iWzqKfdtGyRQt8fHxyXO+xd6GOHjOGqDuJzP/fNh4a8lEIxe0PPM9vmwMYOfIdnJ2dlY5TZjk7
                                OzNy5Dv8tjmA/YHnlY4jRBYWi4X5/9tG1J1ERo8Z89h1H1vQqlWrxpSpU9kfeJ5VGw8WaUghCuP8
                                5VvMWbaVfv364evrq3ScMs/X15d+/foxZ9lWzl++pXQcITKt2niQ/YHnmTJ1KtWqPX5261zHCWre
                                vDmjR4/hl42HWParnwyJJRQXeOIyk75ci7d3M0aNGqV0nHJj1KhReHs3Y9KXawk8IbftCGWZzRaW
                                /erHLxsPMXr0GJo3b57rNnka+K53795MmjSJLX4n+PybjSSnyHA5ouRZLBY27ghm6oL1dHymM59P
                                ny5jNxYhtVrN59On0/GZzkxdsJ6NO4LlUoNQRHKKkc+/2cgWvxNMmjSJ3r3zNsXRI/OhPc6pU6eY
                                /MknqDDxxovt6eLzRLY3twlR1EKvRbJ4lR+nLlxn2LBhDBw4UOlIeVZa50N7nFWrVrFs2TIa16/O
                                iIGd8KzpqnQkUQFYLBZ2+Z/m+98PYEHDZ59/TuPGjfO8fb4KGkB8fDzLly9n06ZN1KvtRv+eLWnX
                                vB5ajXxSFkXvwpVbbN51nF3+p2jUqBEffDCaBg0aKB0rX8piQQM4f/48CxbM58yZM3TxaUzvLk9R
                                v3ZVpWOJcijdZObg0Yus+zuQi1cieOGFF3jrrbewt8/fjNb5Lmj3hYaGsmLFcg4ePIS13oonvWpQ
                                t6YrVSrbY2MtozWIgklLMxGbkMzV69GcOHediKgY6tSuxWsDB9GtW7cyeUagrBY0yPjEvGPHDlb9
                                vJIrV6/h5uLEkw2rU6t6FRztbLCyenQKDyHyIjkljeg78Vy6FsmJs2GkpKbRrl1b3nzzLTw9PQvU
                                ZoEL2n1RUVH4+/tz9OhRQi9d5M7duyQlJRemSVGMzp07h4ODQ6kdlV6ns8Lezo7atevwROPG+Pj4
                                lPkbpstyQXvQ2bNn8ff35/SpU1y5cpn4hASMxjSlY5Uat2/fxsHBQYZfyyNbWxsqOTlRt159mjdv
                                Tvv27alSpUqh2ix0QRNlh8lk4j//+Q/z5s1j5MiRzJ07F71er3Sscq+8FDSRs1u3buHu7s6uXbvo
                                3Lmz0nEqLLnwVYFoNBrmzJnDH3/8wS+//IKPjw+hoaFKxxKizAsJCQGgSZMmCiep2KSgVUC+vr4c
                                PnwYk8lEs2bNWLdundKRhCjTQkJCcHNzw8XFRekoFZoUtAqqfv36HDp0iAEDBvDSSy8xevRo0tLk
                                eogQBXHq1Ck5OisFpKBVYNbW1ixevJjvv/+eZcuW0a1bN8LDw5WOJUSZExISkq/7pUTxkIImGDJk
                                CAcOHODmzZt4e3uzY8cOpSMJUWZYLBZOnz4tBa0UkIImAPD29ubo0aM888wz9OrViylTpmA2m5WO
                                JUSpd/XqVeLi4uSUYykgBU1ksre3Z/Xq1SxatIgZM2bQp08f7t69q3QsIUq1kydPolKpeOKJJ5SO
                                UuFJQROPGD58OPv37+fkyZN4e3tz+PBhpSMJUWqFhIRQo0YNnJyclI5S4UlBE9lq1aoVgYGBeHl5
                                0bFjR+bPn690JCFKJenhWHpIQRM5qlKlClu3bmXq1KmMHTuWwYMHk5iYqHQsIUqVY8eO8dRTTykd
                                QyAFTeRCpVIxfvx4tm/fzvbt22nZsiWnTp1SOpYQpUJqairnz5+XglZKSEETedKlSxcCAwOpVKkS
                                7dq147ffflM6khCKCwkJIS0tTQpaKSEFTeRZ9erV2bt3L++++y4DBgxgxIgRGI0ye7mouI4fP46N
                                jQ3169dXOopACprIJ61WyxdffMG6dev47bffaN++PZcvX1Y6lhCKOH78OE2bNkWjkXnhSgMpaKJA
                                +vXrx+HDh0lNTaVVq1Zs3bpV6UhClLjjx4/L6cZSRAqaKLAGDRoQEBBA3759ee6555gwYQImk0np
                                WEKUCIvFwokTJ6SglSJS0ESh2NjYsGzZMr7//nsWLFhA9+7duXXrltKxhCh2165d4+7du1LQShEp
                                aKJI3B/g+OrVq7Rs2ZIDBw4oHUmIYnX8+HFUKhVNmzZVOoq4RwqaKDLNmjUjODiYtm3b0qlTJ2bO
                                nInFYlE6lhDF4tixY9SpUwdHR0elo4h7pKCJIuXg4MCaNWv48ssvmTx5Mv369SMmJkbpWEIUObl+
                                VvpIQRNFTqVSMXr0aHbs2MHhw4dp3bo1x48fVzqWEEUqMDCQFi1aKB1DPEAKmig2HTt25Pjx49Sq
                                VYu2bduydOlSpSMJUSSio6O5evWqFLRSRgqaKFYuLi5s3bqV8ePHM3LkSIYMGUJSUpLSsYQolKCg
                                IACaN2+ucBLxICloothpNBqmTJnChg0b2LJlC+3bt+fixYtKxxKiwIKCgqhZsyaurq5KRxEPkIIm
                                Skzv3r0JDg5Gr9fTokUL1qxZo3QkIQokKChITjeWQlLQRImqWbMme/bs4Y033uDVV19l9OjRpKWl
                                KR1LiHyRDiGlkxQ0UeL0ej3z58/np59+Yvny5XTp0oWbN28qHUuIPImOjubatWtS0EohKWhCMYMG
                                DSIwMDBz+KBt27YpHUmIXN3vECIFrfSRgiYU5eXlxaFDh+jWrRu9evViwoQJmM1mpWMJkaPAwEBq
                                1qyJi4uL0lHEQ6SgCcXZ2dnxyy+/8N133zF37lx69OhBZGSk0rGEyFZQUBAtW7ZUOobIhhQ0UWoM
                                Hz4cf39/QkNDadmyJQcPHlQ6khCPkB6OpZcUNFGqtGjRgiNHjtCkSRM6duzIzJkzlY4kRKaoqCjp
                                EFKKqSwyHLoohSwWC7NmzWLSpEn4+vryv//9r0yMav79998zb968LBOd3h+c2cnJKfMxjUbDmDFj
                                eOONN0o6oiiEjRs30rdvX27fvk2lSpWUjiMeIgVNlGp+fn689tprODo6smbNmhznnpo5cybvvfce
                                BoOhhBNmdfbsWRo1apSndc+cOYOXl1cxJxJFadKkSaxbt44zZ84oHUVkQ045ilKtU6dOBAYGUqVK
                                Fdq0acOKFSseWWfVqlVMmDCBjz/+WIGEWXl5edG0aVNUKlWO69yfFFKKWdkTEBBA27ZtlY4hciAF
                                TZR61apVw8/Pjw8++IBhw4YxZMgQkpOTATh9+jTDhg1DpVLx9ddf4+/vr3DajNm7NRpNjsu1Wi1D
                                hw4twUSiKJjNZo4cOUKbNm2UjiJyIKccRZmyceNGhg4diqenJz/88AP9+/fn8uXLpKeno9FoqFWr
                                FiEhIdjY2CiW8ebNm1SvXj3H2bpVKhXXrl2jevXqJZxMFEZISAhNmzYlODgYb29vpeOIbMgRmihT
                                fH19OXz4MCaTCR8fn8xiBmAymbh27RrTpk1TNKOHhwc+Pj6o1Y/+eanVanx8fKSYlUGHDh3C1taW
                                Jk2aKB1F5EAKmihz6tevzxtvvEFCQkJmMbsvPT2dWbNmceTIEYXSZRg8eHC219FUKhVDhgxRIJEo
                                rICAAFq1aoVWq1U6isiBnHIUZc6RI0do3759jqP0a7Va6tWrx/Hjx9HpdCWcLsOdO3dwc3N7pOBq
                                NBoiIiJwdnZWJJcouKZNm/Lcc8/JvZGlmByhiTLl7t279O/fP8frU5BxlHbx4kX+7//+rwSTZVW5
                                cmW6deuW5dO8RqOhe/fuUszKoPj4eM6cOSMdQko5KWiizDCbzQwePJjr168/cuTzsPT0dKZPn87J
                                kydLKN2jXn/99SwDLVssFgYPHqxYHlFw96/bSkEr3aSgiTLlww8/ZOTIkZlHOXq9/rHrDx48ONfi
                                V1z69u2b5ZSnlZUVvr6+imQRhRMQEEDNmjWpVq2a0lHEY0hBE2WGWq2ma9eufPvtt0RGRrJv3z5G
                                jBiBq6srkFEwHpSenk5ISAhffvmlEnExGAz4+vpiZWWFVqulb9++2NnZKZJFFE5AQACtW7dWOobI
                                hRQ0USap1Wqefvpp5s+fT3h4OPv372fUqFG4u7sDZB4ZmUwm/vvf/3Lu3DlFcg4aNIj09HRMJhMD
                                Bw5UJIMoHIvFwoEDB3j66aeVjiJyIb0cRYFZLBZu3bpFeHg48fHxj+2oUZKZzp49y549e9i1axdR
                                UVFAxpBUCxcuzPbesOKUnp5O3759sVgsbNiwQbp8Z0OtVmNnZ4e7uztVq1Z97LBhSrh/Q3VQUBDN
                                mzdXOo54DCloIl/MZjOHDh1i586dHDl8mPiEBKUjPVZ8fDxRUVFERUVRrVo1RW5ovn902LBhwxLf
                                d1ljb2dHq9at6dq1K23bti3xDyDZ+fbbb5kwYQJ37tx57JBmQnlS0ESe+fv7s2jhN9wMv8WTXjVp
                                9VQdGtX1wMPNCTtba9Tq0vXJ+mGXroZTt5Z7ie9378ETqFQqOrTNfqaAis5stpCQlMLNiBjOXLrJ
                                keOXOXH2Gh7uVXl31Hv4+Pgomm/gwIHcuXOHrVu3KppD5E4KmsjVjRs3mDdvLkFBR+nYxovX+/rg
                                4eqU+4YCyHjDBkp9wS9NbkbGsPIPf/YGnKVF8+aM+fBDxXoY1qxZkxEjRjBp0iRF9i/yTgqaeKyj
                                R48y5dNPcalsYOSgzjSuL92WRck5deEG3/28m6g7iUyZOrXEr2FdvnwZT09P9u7dS4cOHUp03yL/
                                lD9BLUqtzZs3M378OFo0qcFXkwZIMRMlrnH9anw1aQAtmtRg/PhxbN68uUT3v2/fPvR6Pa1atSrR
                                /YqCkS5XIls7d+5kzpw5vObbloG+7UpdzzNRceistPz77WfxcHNizpw52NjY0LVr1xLZ9759+2jV
                                qhXW1tYlsj9ROFLQxCPOnTvH7Nmz6NujBYP6KHtBXgjImKVgUB8fkpONzJo1k6pVq9K4ceNi3+/e
                                vXt58cUXi30/omjIKUeRRVxcHBMnTMC7UU3efLmj0nGEyOLNV56h2RO1mPzJJ8TFxRXrvqKiorhw
                                4YJcOytDpKCJLFasWIEKE/8Z/qz0yhOljlqt4j/Dn0WFiRUrVhTrvvbs2ZM5IasoG6SgiUxXrlxh
                                8+bN/Oulp7GxVmYeMSFyY2Ot440X27Np0yYuXbpUbPvZt28fTz75JI6OjsW2D1G0pKCJTEuXLqFu
                                LVc6t2ukdBQhHquLzxPUq+3GihXLi20f0lW/7JGCJoCM6wWHDgXQv2cL6dEoSj2VSkX/ni05dCgg
                                c7zOonT79m1OnDhRYr0pRdGQgiYAOHDgANZ6K9o1r690FCHypF3zeljrrfD39y/ytnfs2IFKpaJj
                                R+kYVZZIQRMABAcf5UmvGmg18l9ClA1ajZonvWpw9GhQkbe9c+dOWrdujZOTDPFWlsi7lwAg9NIl
                                6tZ0VTqGEPlSt6YrocXQMWTnzp1yurEMkhurBQDRt2/j4pzP0eBNsRz8dQ0L1gcScD6K28lg6+xC
                                gye86NypNc92aUar6jYV8FNTEMObTOHnpIceVmmxtTfg4lGT5i288X21J/2bOBbLH6EpNoydm3fx
                                y5+B+J++SUSKFc5VXaj/VFuGjfClfyP7fL4uaVzxW8+EaavZFNaaXy+M44UC5Io/vYNpX27ij6Dr
                                3DbZUse7NcM+HMrIFg4U5Mptlcr23Ll7twBb5uzq1auEhoZKQSuDKt57jchWaqoRvc4qH1vEsH7s
                                GHp8ehBL97fYuO17boWs5NjqDxn1xG1WTp1Blz4rOFBsifMp6Rgfdn6FOm9tI6zYd9aCJSGbSNz8
                                Cp4A3ccSG7qJxPO/cWnHLL7/sB1uoX8x3HcY7cbt5Fxy0SfYNeM/9JtyAHO3t9i47SduHV3Kjlkv
                                UP3MOob6jmbkttg8t5V87TCfvz2KzrMCCY1OLXCmxGM/0+Olr/nTrhe/bl/FjX2fMa7mWSYN+Dcf
                                7osvUJvWeiuSk1MKnCk727dvx9bWlnbt2hVpu6L4SUETQMZMz/nq3XjiDyZuiqbKK+/x44hWNKlq
                                j43eBreaDXnl35+wYpBb8YUtCIsZs8WCxWJGseklNDocqnjQuvsLfPXzAjaPqEHo2nn0Gr2Ta8UQ
                                yuWVd1nxL28authibWNPnVY9WDyvH/VNUfz8xQZC8tSKiY1fLeN0i1H4b3qLTnYFDGO+wlfjf+OE
                                QzcWznqWFq56bCrV5tXPPmJUrXCWTljJfmP+m1WpVEU+U/rOnTvp0KEDer2+SNsVxU8KmiiQhPPX
                                CANqe1bj0Vuw9XR8vhUuJR8rZ4bmzPdbw5UVvaipdBYA7Hl63Dgme2uI3LGECZsLdoSSk+5f/MqV
                                6d48PL+yplFDWlgD18K5lKc6oMF35jesGvkU7oWYrNl0eCsrLlhwf7YznR4c51fjycAXakH4Lpbt
                                KvjRX1GxWCzs3r1bTjeWUVLQRIHYVXHCAJzeH8yt7FZoM4IrQaOQ21IfQ1WVt4a0xIokNv60m4iS
                                2GdSPHdSgAa1eCKPB+RFMWrMmYMniAKaNa33yDKvJ+thSwq7D5wr9H4K68SJE0REREhBK6OkU4go
                                mNbt6euynZ/3fkfXf91k/Ihe9G/tjl0uH5Esd86x7OvVfL/jDOciklHZV+GJ1m0Z/v4ABj1hyFhp
                                21c4jPTDBEArvtvxLFfn/Moq/ytcjzHeezyrpmO/5dB71cG0l1frzyZz1qxe/yGx75Es7S0/818G
                                PHA2yXL3Aj8sXMP/dpzidHgyWmdXPD0b0KtvD97o3YQaDxxR5Cl/Pti3bEwTAggOPoV/ui/9svuL
                                jNvG895f45fX57so53unbv+5n73Y0X/Uc5TcHYcWzl26CdhRzf3RaVhUbs5UBUIv3wCeLLFU2dm5
                                cyfOzs54e3srmkMUjByhiYKxbcGshUPo6g5X9qzjnYHDqdnmPXp/uJwFm04Tll1Hh0h/3uwznjFb
                                Uujz2SzOHlvFyV/fplvMboa/OIGpR++dcurxEXGhf7C8O8AlZkzahevr/2b/gV+4tG4AzTVNmB/0
                                M4ue0YG6LtP9NmW8uQNoOvJb6ApmNDfQZ/4vGW/uWdp7SFQAb/cdx/ubknnuv//HiaO/cmHjx/y3
                                bSzLx03kX6ui858/P1yccAMw3eXWnRzWcejBltB8PN+cRPszelYQ7q9+xILnHfKftcCSiY0zAdbY
                                2mSz2GCDASAuoQQzZW/nzp106dIFtVreGssiedVEgTm1fJmNu5ey/cvBDOvekKopYeze8AcTR4+n
                                8dNjGbPpBv9c5zfy1+yFrL6hw/eT8YzrVJ0qtjZUrd+ayQuG0MlyhS+nbCL0kb2k8cyIMQxv605l
                                Gx0u3oPYd2EGwyo5MODtrriZL7Fo2QnSHtgiPWgDC291YfRzufVgMLJ19tf8EmZF3/9OYHyXWrgb
                                dNhVqUnP9/7D1Ge0WdYtWP5cWMhjJ5VCPt+7J5jwxjyC2n7In5+3pFJ+cxYni+Xev4GyQ66lp6ez
                                b98+Od1YhklBE4Wjq4JP/1eYv/hLTh/7meCfxzDxhdoY7l5g6djpzDx1/wThaTZtiwN1I57v8tDR
                                gctTdG4A5pBD/PnIBbnatHwq+95mep++jGys4sbv61mbeStSImuX7KD6m31ok2snhtNs/DsW8KJH
                                p4dPF9ox9H/r2fFmlULmz0XknYxrkNrKuFcGzq/C2/MFDA981fn0ZOGeb9J5pg79nI313mf73E7U
                                LETnjoKxwdFBA6SQlN2Re1IKSQAO+T9lW5T2799PfHw83btndygvygK5hiaKjsaOBu268km7Z3iu
                                2gd0+C6M3/+6yuTGnmCMITIe4CgjnnyBEdk2cJNLl4GqDz5mjcE2px168PawNnz5YQDzV4bx2vs1
                                4PJfzD/cjIlz8nDbwP1MeidccnsvLXD+x4sLPM0pQNW8Me21QIOBHAsdmMPaBXi+pnCWvzeNH6u+
                                zfavOlC9xIsZgIqGdT2AMG6EpwBZr6NZIm5zC6hSp5oS4TJt2bKFRo0a4enpqWgOUXByhCYKJmgZ
                                dVt9y/5sF2pp3rYRdsDd2HvXRXSVcHMAND6svLCJxNDsvlYxN5/3slbq3Z/B7hZO/rieXcZ0di3b
                                ROKA/vTOy4d9nROu9kBqDFGJua1bDPnN4Sz7KYh0DPQb3ClPtznk7/nG8ufHn/JZ6sts+qY7npnF
                                LJTJnd5gYnA+shZSo7ZNcQGOhTw6TNW5k5dIQk/n9g1LLlA2Nm/ezPPPP69oBlE4UtBEwVgsmG8H
                                syU4u7thLVw4GUoCap5qXOveY43o09MRTOfZH/joNhe/G4t9+yUcyq4L4+NoGvH+v7zQ3PZj3rL1
                                zNvkwftD6+fxP/YT+PZ0BM7y9+6HK1o4s3r70ujzM5iLJX88+2fPZvpxE1V7jmDGc/Z52yzPzzeV
                                wPnTGHWhE78v7cMTCs/XqmnzLG/WV3HzLz/2Pth3xnSFXzZfAfeuvNVZuRuZQ0NDOXv2rBS0Mk4K
                                miiEcL557zOmbTzDxYhEUowpRF+/xOals3l5wUUMTV5lSv/7s/3q6DHuA16vGcOy8bNZ4HeN8Hgj
                                yTHh7F81j1e+vkv/j/vTtgCnxGq/1g9f+zR2fvUTx3r04/U8n/LT0es/7/NaDSMbPv+CWbuvcSvR
                                SPytC/w2eSazI1sx5c1G9/5IiiC/OY2E2+Ec2b6Zfw/6gBcWh1H35Q/ZOrcz1fPRHyL352vhytqZ
                                vDj/PJHHfqFjk6zX5Ayeo5lz7aFNYo/wXpdXcOs4izWFuSEup3bUtfnoi1d4MmYb747bytGoVFLu
                                XmX1p1/x9WV3hs14nQ4KDsyxefNmHB0dad++vXIhRKHJNTRRMM0Hs3tNIzZtC2D794tYMyOa8Ohk
                                zNb2VPesS8cxHzPmjXbUe/DIwLk13/0xG5+Fq1kxZQKf3UxG7VCZOo2b868ls3in/b0OGMErqP/i
                                em4CcIS3Gr3AW/iwMnQi/bLLYmjHmNfcWb9EzfBhrXikZ3iW+9oy2nu3zyTuzG0LLm1Y+sds2i9c
                                zf+mTGBWeArayi40bvMM/1v9Er09CpD/4cGJt8/B0XMOqDTY2Nnh4lGD5i2fZcnEnvQryODEuT1f
                                0tj/1xGis9k0RxYzZrMFi9nyaK/LXfOoPGwn/xxY7WOA5z4Aus34lQ2vGvLUjqHZ6/y91o1pX23i
                                la5LuGO2pfZTrfn816G80yKPR6jFZMuWLfTs2RMrq/yMZypKG5WlqAdCE2VS586dmfDOC3Ro1UDp
                                KELk2b4j5/ni203s3r27wG0kJibi7OzMkiVLGDJkSBGmEyVNTjkKISq0bdu2kZaWRq9evZSOIgpJ
                                CpoQokLbsmULrVu3xtVVJrgt66SgCQA0Gg0ms1npGELki8l7Azk6AAAcJ0lEQVRsRqMp+M11FouF
                                rVu3Su/GckIKmgDAYGtLUpLy03cIkR+JSSkYbHO88z5XR48e5caNG/Tu3bsIUwmlSEETALi7V+Vm
                                RNFOZS9EcbsZEYOHh3uBt9+yZQvu7u489dRTRZhKKEUKmgCgfoOGnA0tkRm5hCgy50IjqFe/4D1z
                                //zzT3r37p2/2dpFqSUFTQDQunVrzoXeICYuKfeVhSgF7sYmcvbSDVq3bl2g7cPCwjh8+DC+vr5F
                                nEwoRQqaADIKmsHWwLZ9J5WOIkSebN8fgp3BUOCCtnbtWhwdHWV0/XJECpoAQK/X0/uFF/hjezAJ
                                iSlKxxHisRISU/hjezC9X3gBvb5gY2atWbMGX1/fAm8vSh8paCLT66+/jlarZ9XGg0pHEeKxft5w
                                EJVKy6BBgwq0/fXr1zl06BAvv/xyEScTSpKCJjLZ2try5ltvsXnXMS5djVQ6jhDZunQ1ki27j/H2
                                8OEYDAWbFHTt2rXY29vTrVu3Ik4nlCQFTWTRq1cvvJ/yZtrXG7gTk9skYUKUrDsxiUz7egPeT3kX
                                aqiqNWvW0KdPH6ytrXNfWZQZUtBEFmq1milTp2JrcGTa1xtISU1TOpIQAKSkpjHt6w3YGhyYMnUq
                                anXB3r5u3LghpxvLKSlo4hF2dnbM+OILou4kMWHWGjlSE4q7E5PIhFlriLqTxIwvZmJnZ1fgttau
                                XYvBYJDejeWQFDSRLQ8PDxYuWkRqupYPP18l19SEYi5djeTDz1eRmq5h4aJFeHh45L7RY8jpxvJL
                                5kMTj5WQkMDUKVMIPhbM8529GdSnHXYGeSMQxS8hMYWfNxxky+5jNPNuxqdTphTqyAwgPDyc6tWr
                                s379ermhuhySgiZyZTab2bp1K0uXLMFiSadv92Z0f7oJlRwL1sNMiMe5G5vI9v0h/LE9GJVKy9vD
                                h9OrV68CXzN70IIFC/jkk0+IjIyUI7RySAqayLPExER+/vlnNm/aREJiIg09PWhUtyrubk7Y2Vqj
                                Vst4eCL/zGYLCUmphEfc5cylW5wLvYmdwUDvF15g0KBBBe6an50OHTpQq1YtVq5cWWRtitJDCprI
                                t9TUVA4fPsyRI0c4f+4s4eG3SEhMxCzzqRWL2NhYYmNjqVmzptJRioVKpcLOYMDDw536DRrSunVr
                                2rRpg06nK9L9XLt2jTp16rBu3Tr69OlTpG2L0kEKmhCl3I8//si//vUv/Pz86NChg9Jxyqxp06ax
                                cOFCrl+/jpWVldJxRDGQgiZEGdCnTx/Onj3LsWPHsLGxUTpOmWOxWKhfvz59+vThq6++UjqOKCbS
                                bV+IMmDRokVERkYybdo0paOUSX5+fly6dImhQ4cqHUUUIzlCE6KM+Pbbb/nggw84ePAgLVu2VDpO
                                mTJkyBDOnTtHQECA0lFEMZKCJkQZYbFY6NGjB5GRkQQGBsp1oDyKi4vD3d2dr776ipEjRyodRxQj
                                OeUoRBmhUqlYsmQJly5dYubMmUrHKTN+/fVXzGYzr776qtJRRDGTIzQhypivvvqKjz/+mKNHj9K4
                                cWOl45R67dq1o27dunLvWQUgBU2IMsZsNtOhQwdMJhMHDhxAo9EoHanUOnfuHI0aNWL79u107dpV
                                6TiimMkpRyHKGLVazfLlyzl+/Djz589XOk6ptnjxYmrVqkXnzp2VjiJKgBQ0IcogLy8vJk6cyOTJ
                                k7l48aLScUqlpKQkfvjhB955550iGQdSlH5yylGIMio9PZ3WrVvj4ODA7t27UalkLM0Hfffdd4wd
                                O5awsDCcnZ2VjiNKgHxsEaKM0mq1rFixAn9/f5YsWaJ0nFJn0aJFDBw4UIpZBSIFTYgyzNvbm7Fj
                                xzJu3DjCwsKUjlNq7Nq1i5MnT/LOO+8oHUWUIDnlKEQZl5qaSvPmzalVqxZ//vmn0nFKhX79+hEd
                                Hc2+ffuUjiJKkByhCVHG6fV6li1bxt9//y33WpExTcymTZt4//33lY4iSpgUNCHKgXbt2vHuu+8y
                                ZswYIiIilI6jqIULF+Lm5ka/fv2UjiJKmBQ0IcqJL774AkdHxwp9ZJKcnMzy5ct55513ZKzLCkgK
                                mhDlhMFgYOnSpaxdu5bff/9d6TiKWLlyJQkJCbz99ttKRxEKkE4hQpQzb775Jn/99RenTp2icuXK
                                SscpMWazmSZNmuDj48OyZcuUjiMUIEdoQpQzc+fORavVMnbsWKWjlKg//viDs2fPVrjnLf4hR2hC
                                lEObN2/G19eXP//8k169eikdp0T4+Pjg5ubG+vXrlY4iFCIFTYhyasCAARw8eJCQkBDs7e2VjlOs
                                9uzZQ6dOnThw4AA+Pj5KxxEKkYImRDkVHR1N48aNeemll1i4cKHScYrV888/T3x8PHv37lU6ilCQ
                                XEMTopyqUqUKc+fO5dtvv2Xnzp1Kxyk2J0+e5K+//mLcuHFKRxEKkyM0Icq5fv36cfz4cU6ePInB
                                YFA6TpEbPHgwQUFBhISEyDQxFZy8+kKUcwsXLiQmJoYpU6YoHaXIXb9+nd9++40JEyZIMRNS0IQo
                                7zw8PJg5cyZz5szB399f6ThFavbs2bi5uTFgwAClo4hSQE45ClEBWCwWevXqRVhYGMHBwej1eqUj
                                FVpERASenp7MmDGDDz74QOk4ohSQIzQhKgCVSsXixYsJCwtjxowZSscpErNnz8bBwUGGuRKZpKAJ
                                UUHUrl2bzz//nOnTpxMcHKx0nEKJjo5m8eLFjBs3DhsbG6XjiFJCTjkKUYGYzWaeeeYZEhMTCQgI
                                KLMj0k+YMIHly5dz+fJl7OzslI4jSgk5QhOiAlGr1SxbtowzZ84wb948peMUyO3bt1m0aBHjxo2T
                                YiaykIImRAXTsGFDJk+ezOTJkzlz5kyWZQcOHGDQoEEKJcvKYrGwYMECkpOTszw+d+5cdDodI0eO
                                VCiZKK3klKMQFVB6ejpt27ZFp9Oxf/9+EhMTGT9+PN999x0qlYrbt2/j5OSkaMbw8HA8PDxwdXVl
                                2rRpvPnmmyQlJVG7dm3GjRvHxIkTFc0nSh8paEJUUMHBwbRp04YRI0awfv16IiMjSUtLA2DTpk30
                                7t1b0Xz3BxyGjFOlNWrUoEWLFvj5+XHlypVyP+CyyD855ShEBeXp6YmPjw/ffPMN4eHhmcVMp9OV
                                ikF+z507h0ajATI6s4SFhbF+/Xqsra3ZunUr8llcPEwKmhAV0JYtW2jYsGHmyCFmszlzmdFoZNu2
                                bUpFy3ThwoXMggYZGS0WC7du3eKVV16hZcuW7Nq1S8GEorSRgiZEBRIZGcnLL79M7969iYqKyjwq
                                e1hISAhxcXElnC6rc+fOZZvvfvE9fvw4Xbt2ZdiwYSUdTZRSUtCEqEBOnDjB9u3b0Wg0WY7KHmYy
                                mThw4EAJJnvUqVOnHntaUa1WU6dOHSZPnlyCqURpJgVNiAqkW7dunDx5kmbNmqHVanNcT+nraCaT
                                ibCwsByXW1lZUadOHfbv30+tWrVKMJkozaSgCVHB1KhRg/379/POO+8AGeM8PsxoNLJ9+/aSjpbp
                                ypUrOZ4O1Wq1NG3aFH9/fzw8PEo4mSjNpKAJUQHp9XoWLFjAypUr0ev12R6tHTt2jPj4eAXSZXQI
                                yY5Wq6V9+/bs2bMHZ2fnEk4lSjspaEJUYIMGDSI4OBhPT89HxnU0mUwcOnRIkVznz59/JI9Go+HZ
                                Z5/l77//liGvRLakoAlRwXl5eREUFETfvn2znH7U6XTs2bNHkUwXLlzIkkWtVjNw4EDWrVtXLuZy
                                E8VDCpoQAjs7O1avXs13332HVqtFo9FgNBrZsWOHInlOnz6N0WgEMq7xjRo1ih9++OGxHVmEkKGv
                                hBBZ+Pv78+KLL3Lr1i20Wi1xcXElPudYtWrVuHnzJgAzZsxgwoQJJbp/UTZJQROiAjKbzSQmJmb+
                                npKSQnp6euaysLAwxowZQ0BAAEuXLqVZs2aPbS89Pf2RUfEfx8rKCmtr62yXGY1G2rdvD8CUKVMY
                                MGAAkNGRRafTZa734HW0h5eJikkKmhClzP1ik5CQQEpKCsnJySQnJxMfH09ycnLmYwkJCSQlJWEy
                                mUhISMBisRAXnwBA3L3eifd7KSYmJNz7nr9ei9euXcNsNlO7du2ie4K5SExMJDAwkEaNGuHq6prv
                                7fXWNmi1WvTW1mi1WmxsMn63sbbGysoKGxtrdPcKqlarxc7ODhsbm8yvnH43GAxldkLUikIKmhAl
                                ICYmhqioqMyv2NhY4uLiiI+PJyYmltj4eOLi4oiLjX1s0dHqrNHqM76sbO1RW1mj1lqh1tuiUqvR
                                WhsAFVa2doAKrfU/jz+4HLi3Tga1lQ6N1T+dLbQ2dpmdMuJvR1GpWu1cn+P9feRFekoSFrMp22UX
                                Du0gPc2IV/uemIwpmY+bUpMxmzKOIi1mE+kpSQ8sS8lcZjImZy7P+H7v99QkLCZT5nJz6r3Hk+Ix
                                GZMxGVNJS03GmJSQY2693ho7BwccHBxwcnTEyTHj5/tfzs7OuLi44ObmRuXKlbOMRSmKnxQ0IYpI
                                WloaoaGhXLhwgfPnz3Mt7DqRUVHcjorCaEzNXE9vcEBv74jWxh4rgwMaG3usbB3QGRywsrVHZ3BE
                                a2NAa2OHRmeNVm+T8f2BIlOeWSwWxZ+nKTWZ9HtFLj05gfSUJEypyRgTY0lLiictKT7j54RYTCkJ
                                pCfFYUyIIyU+BvO9Qq1Wq3GsVBk3NzequVelXr161K9fnwYNGsjUN8VECpoQheTn58fqNWs5f/4c
                                pvR0dNa2OFSri7VrDawdq2BTyQ1rpyroHV2wruSS5UhIlC8Ws4nUuLuk3I0gJTaalJgoUmKiSL4d
                                TsKNiyTFRAPgWtWdHt26MmTIEDmNWYSkoAlRCMHBwYz96COqNGiOR8vuONRogG0VjzyfehMVizH+
                                LrHXLxB9OoDwoJ307NaFjz76SOlY5Ybc1CFEIcTExIDFQvS5IO5ePoV9tbo4Vq+PfbV6GFyrY1PJ
                                DZ1DJVQqKXAVUVpSHKmxt0mKvkncjUvEXb9A3PWLpMbdRqPTExkZqXTEckUKmhCFcP9az9MTlhNz
                                +RRx1y8Qc/Us1w/9hSkt47qZSq1B71AZ60quWDtWwdqxClZ2jugMjlgZ7l87c8DKkHEtTa2RP8vS
                                LC0pjrTEe9fQkuJIS0rI+Dkx7t5pxmhSY6NJuRuZ+X8AlQqDSzUcqten9jP9cajRgDD/Tdjayged
                                oiR/OUIUAYNrDQyuNajWphdw/1rKHVLuRpJ8N5LU2GiSY6JIiYnk7uVTGBNiMCbGYUpNeqQtrd4W
                                K4M9Wr0tmnsdQqxsMzqIaHTZ/66x0qG20qHWWqHRWaPWWKHRW6PWZjz2cC/G8u7BXo4mYwoW071e
                                jvceN5vSMaWmYDGlY0pLIS0pHlNqCiZjCumpyaSnJGBKzfjZZLy33JiCMTEWHrpKo7HSZ7xeNvbY
                                OLlg61yVyp5NsK7kht7RGWsnF2wquz3y7x/mv6kk/0kqBCloQhQDlVqDtZML1k4uONVpnON6ZlM6
                                6Q/2mkuMv/epP4701JSM3nb33mCN8TGkp4Zn/J6ckPnma05LzbH9h6m1OjQ6fWZGje6fm5szuvhr
                                MpdprW3/Waa3QZWHLuhqKz1qzWM6OaggPTnnbvEPSk9OzJzg05yelqULf3pKIpZ7E5SaTWmY04z3
                                fjZl+yEhxzj3nqfWxpDRk1Rng0Zvg5WtHTp7J2yc3TM/QGR80HjwiDrjZ3UF+qBQ2klBE0JBao0W
                                nX0ldPaVMBSiHVNaKuY0I+b0NMzpRkzG1HvfUzIeSzNiSjNiTkvFbMqYZyzjsX+KYW4FBHLvP5Ya
                                ezvXdTQ6PWpt7j379I5VMgusWqNFo/9n+C2ttU0OxVf1wD13ud+bJ8oXKWhClAMaK32FOqUoRHbk
                                iqQQQohyQQqaEEKIckEKmhBCiHJBCpoQQohyQTqFCKEEcyRX181i3+a/uHYpjMQU0FWqgUvDttR7
                                ujdeHbpRw8Oe8j8U8cP+Zk37Fzj68NRqKius7CphV7UR1by70KTvWzRt5JrHT+RF0Ka8XmWCHKEJ
                                UeIiCPmkLYu/2ICl8yzeXHuFKQduMHbFCtp73SRo1qt8O2gCl5WOeV/SDjb0cWb66BXEFPvOevLy
                                ASMzfpmIM0Cn75l+1MiMI7f5eJ0fA97pi/2VJawZ1JAFU34kMiW39oqizTL2elVgUtCEKGmn5rNl
                                63UMfRcxcOjzVHWrjJXOHvvqbXhq1O+8+nJtpRNmZTFjMVvg3o3MilBbY+1cj5qdRuG7OIi3hnpx
                                e+Mwlk78kZiCDq+e1zbL2utVgUlBE6KEGS+dJgaoXKshj469YYtn9+ewe3Qz5Rh60HfTHSZ9PQwn
                                pbMAUJk6H6yiexMtCXvGsvnv3G/mLkybZe71qsCkoAlRwnTOruiAW4e2k+3c1C3mMWn3QjxLOFeZ
                                oqpDmwHPoiGO06tXZf/vWERtyutVdkinECFKWrMXaeL8PUcPjubb9y7S5Y23adrcE30uHy8tdwMI
                                WPoFR/YcJCoqAeyqUbVZH9oOn0Tzho4ZK+0eyqSPfiHj5OBzvLRuOHe/nc7RwyHExiaT3UlD93dP
                                8sGwhmBazU+tXuf0/QVdVzLjuS1Z2nvl0B800z2QKTaIwGVfcMTvALci49FUqknl2q3weu4tWvXs
                                gNMDg5fkKX8+6L3bU5VN3Di5n6vp79OkCN7Nsm2zgK+XKHnykghR0mx70nv2Z9R3g7v+X/H7cC8+
                                79GM5Z+MY9/WA8Rk19Eh+g9Wv96ZDduTaDJxD+P9wvnP0jnUj/uZNUO7sO3EvQF5O//A9KNJvNIJ
                                IJid01di9/JPvPdXBBN/nER1dUf67rrJiz42oG7GsxuNGcUMQPMKg49e4rknnWj8fxHMmP3KQ+09
                                5PYm1rzekfVbE/Aat51/74pk4s9r6NEymsNTuvLr7zfynz8/qlTFHsB0i7ii6q2SXZsFeb2EIqSg
                                CaEAG+/xvLnhHCOmTaVNp9bYp57l4p/z+PPjzsx+1ocNf1/AlLl2Mme/fpdj4TY0HvsznZ9uiMHW
                                Hvu6z9P9i8+oy0n8Zn7Do1eSjNR9YxltW3pia22NXZNPGRW4gzZOVfAeMhh7czD+K/0e2A+Yj83n
                                QOTrdOxeKZdnkMzZr0cSfMOaJuN+oUuHxjgYbNA5P0HDYT/R08fqoXULkj8XFksehksumjbz93oJ
                                pUhBE0IpumrU7j2RvnP2M87vJmMXL6Nrz6boYgM5NOkldp1Nv7eiP6d2R4O6HY06VsnahnMX6nmC
                                5cxGzkY8vIMmVG9i+/CDAGhbj6Gdl5rYjXM5EXv/0RiO//gjToM+oGauM8X4c2pXFNCWhj4PdxWp
                                RMtvEhk5sFoh8+ciOjzjmpbWHQcn4OJnzGmuY+IDX9O/2FO4Nh+U59dLKEWuoQlRGmgq4dJqCN1a
                                DcDLvSULvz/DiZ0hdPfyBmMkCQkA21jbQcfabBu4SPQ1wO3BxwzobLJdGahHm9d74/fJJvatPkOz
                                txvB1SXsO9qNrp/Vzj3v/Uw6Fwy5zcJS4PyPlxJ8gFuA6skO1NYC9SYz9ujkvDeQlzZz8rjXSyhG
                                CpoQJe3Yv/m/fxt5bccC6jyyUEf1Vj7ovj9Lctzdew+5YW8PJPVj0KHfaJL7PJt5YtvzI1p8vZGD
                                v83lwtBv4KeFGPv9zhN5mSZM54qdHZAQRWIiPHYyt+LIbw4lYPXfmHGk6auvFU23+ZzazO/rJRQj
                                pxyFKHEWLHe2ceZkdr0JLESfPoYRDR4Nm9x7zIfGXVzBdJjQY49uE/29Dx8/+yFX83sRR9OODgPb
                                or7zC/t+msu+v+vx9Gst8zh8kw+Nu7gAhzh74OEeGZfY/ZqemV/637seVdT573D560HsCEnHvvM8
                                nu9WOa8bFrDN/L5eQilS0IRQRCj7x/Vn+1/+REfFkG5MIPFmMKd/fJ0flxxF1+hjevR2ubeuNQ3f
                                X0yL6pEETHmd/QdOE5eQQlpsKJd/f4ufloTTdOy/qVWAI59K/T/kCbtULiz6Lzc7j6WFa163tMHr
                                /e9oVi2ZU1++xu59p4lPTCY1IohjMwayO/o5er7uc684FkF+cyqpdy4R5reITSNasvyHs1Tps4K3
                                /28QjgUdQDFfbebn9RJKUVnuz7kuhMg3Pz8/pk6dSs852/K+kTmBuyf+5pTfRs4fPcmdiBvE3YnH
                                Yu2MY61meHZ5i44D+1DFOutmltijBC2fweHd+7kVEY/KvirOXj1oPmQiPm1qZHw6PTmeGUPnEpdl
                                y34MOvob2R8/mAmb35hFP2jo+ttJutV/6J08y31tGbTPreWzz33vZQokcOkXHN5zgIiIBNRONaja
                                4jU6vvNvnqiZ9QJenvLnOJCwFitDJQxVvajerBtN+g2jqZdLIQcnzmObBXy9cnP8x8/xclLz6aef
                                5m9DkSMpaEIUQoEKmhBIQSsOcspRCCFEuSAFTQghRLkgBU0IIUS5IAVNCCFEuSAFTQghRLkgBU0I
                                IUS5IAVNCCFEuSAFTQghRLkgBU0IIUS5IAVNCCFEuSAFTQghRLkgBU0IIUS5IAVNCCFEuSAFTYgi
                                YEpLVTqCKEMsFgsmo1HpGOWOVukAQpRlzZs3p1JlZ/ZOGUCluk9iX60e9tXr41ijAXqHophJWZQH
                                ptRk4m5eIv76RWKvXyDucggJ0eF0HzRd6WjlisyHJkQhxcXFsWPHDk6cPMm58xeICL+JxWJBZ22L
                                jVMVdE6u6ByrYOPkgnUlV3T2ldDbV0Zn54iVrT0aXT5nhhSlhsWUjjExlrTEeNKS4kiJvU1qbBQp
                                MdGk3I3AGBtNStxtkmPvAGBn70CDBg1o5NWQrl27UqdOHYWfQfkiBU2IIpaYmMiFCxe4ceMGUVFR
                                REREEBkVTWRUFJG3bmE0Zj09qdFaoTc4oLNzQGvrgMbGHitbB7TWtmj0Nmis9FjZ2qHR2aDR26DV
                                W6O1NqDVZyxXa63Q6KxRa60UesZlj8ViIT05AYvZhMmYQlpSAqa0FEypyaSnJJOekojJmILJ+M/v
                                6SmJpCfFY0qKw5gYizExHmNKUpZ2NVotTpUq4+rqikdVN6pUqYKrqytVq1albt26uLm5KfSMKwYp
                                aEKUsKSkJO7cuUN8fDxxcXHExsYSFxeX5SsmNo7EpCSSk5MzvpKSSE5KxGw2P7ZtrU6PxkqHxkqP
                                WqtFq7NGpbVCY5XxXaXVobbSAaDWWKG20v+zrY0dKpUqY5nWCo3ugWXWBlTqXC65WzLayI3JmILF
                                nJ7reukpSVjMJgDMpnRMqSn/LEtNwmLKWJZRlJLv/WzGnJqExWzGlJoIFgtpyQkZ31OSsFgsGJMS
                                ct23tY0t1jY22NjYYGtri53BgMFgi6ODA46Ojplf9vb2ONx77P53dW7/TqLYSEETogxJTU0lJSWF
                                pKQkEhISSE5OJi0tjeTkZNLT00lOTsZkMmX5/eHvpnuFIDk5BWNaGpBxxBKfkJBlP8YHOi0kJuRe
                                BNLSjBhTc+8co9ZosLGxzXU9vbU1Wm3GZX61Wo3BYMhcZmNtjZVVxhGpSqXCwf6fQmpnl/Gzvb19
                                lu92dhkF+/5yOzs7NBoNBoMBW1tbbB4oYKJskoImhBCiXJBjYyGEEOWCFDQhhBDlghQ0IYQQ5YIW
                                WKN0CCGEEKKw/h+IhFKBUi5eCAAAAABJRU5ErkJggg==" alt="example">

Acknowledgement
-

A thank you to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>, for letting me share this work under an
open source license.


