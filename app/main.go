package main

import (
	"fmt"
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/opts"
	"github.com/TopoSimplify/offset"
	"github.com/TopoSimplify/constdp"
)

var constraintsWkt []string
var linearWkt []string

func init() {
	constraintsWkt = []string{
		"POLYGON (( 486.7703509316769 722.8914456521738, 475.4937795031055 708.0909456521739, 495.9325652173912 703.8622313664596, 505.7995652173912 713.7292313664595, 519.8952795031055 709.5005170807453, 529.7622795031054 725.005802795031, 510.02827950310547 729.939302795031, 489.58949378881977 729.939302795031, 486.7703509316769 722.8914456521738 ))",
		"POLYGON (( 543.8579937888197 460.00637422360245, 516.371350931677 426.8814456521739, 537.5149223602483 426.8814456521739, 537.5149223602483 440.27237422360247, 562.887208074534 434.6340885093167, 581.2116366459626 454.36808850931675, 562.1824223602483 462.8255170807453, 543.8579937888197 460.00637422360245 ))",
		"POLYGON (( 556.1038053880469 503.83863973683697, 551.5530679870658 503.6964291930563, 550.4153836368205 499.5723234234172, 552.264120705969 496.44369146024263, 555.3927526691436 491.60853297170024, 559.2324373512214 493.5994805846295, 556.3882264756082 499.71453396719784, 556.1038053880469 503.83863973683697 ))",
		"POLYGON (( 266.6254208102681 549.9800181165591, 250.83622260496082 501.7352458225647, 277.15155294713963 490.3319360076206, 299.9581725770279 551.7343734727044, 266.6254208102681 549.9800181165591 ))",
		"POLYGON (( 423.9302472540722 489.78849755284546, 423.9302472540722 512.734076819228, 444.4605023871514 512.734076819228, 444.4605023871514 489.78849755284546, 423.9302472540722 489.78849755284546 ))",
		"POLYGON (( 273.9909412231634 433.60761308851255, 274.819826756752 455.15863696181606, 292.77901331783823 455.15863696181606, 291.39753742852395 432.5024323770611, 273.9909412231634 433.60761308851255 ))",
		"POLYGON (( 219.83535798993222 660.6726799840633, 219.83535798993222 678.1837799505131, 235.5349648564045 678.1837799505131, 235.5349648564045 660.6726799840633, 219.83535798993222 660.6726799840633 ))",
		"POLYGON (( 381.05824388793627 358.75716332113427, 381.05824388793627 381.098911554191, 398.56934385438615 381.098911554191, 398.56934385438615 358.75716332113427, 381.05824388793627 358.75716332113427 ))",
	}
	linearWkt = []string{
		"LINESTRING ( 533.1695652173912 534.3086956521738, 509.0195652173912 551.5586956521738, 473.3695652173912 528.5586956521738, 463.0195652173912 560.7586956521739, 430.8195652173912 571.1086956521738, 484.8695652173912 618.2586956521739, 523.9695652173912 594.1086956521738, 572.2695652173912 613.6586956521738, 581.4695652173912 497.50869565217386, 453.8195652173912 475.65869565217383, 406.6695652173912 506.70869565217384, 372.1695652173912 504.40869565217383, 360.6695652173912 461.8586956521738, 383.6695652173912 422.75869565217386, 427.3695652173912 410.1086956521738, 507.8695652173912 359.50869565217386, 632.0695652173912 354.90869565217383, 628.6195652173911 421.6086956521738, 624.0195652173912 488.3086956521738, 681.5195652173912 560.7586956521739, 724.0695652173912 712.5586956521738, 617.1195652173911 680.3586956521738, 449.2195652173912 683.8086956521738, 333.0695652173912 774.6586956521738, 344.5695652173912 674.6086956521738, 334.2195652173912 545.8086956521738, 231.86956521739125 446.90869565217383, 144.46956521739125 418.15869565217383, 337.6695652173912 308.90869565217383, 489.4695652173912 239.90869565217386, 735.5695652173912 212.30869565217387, 796.8760652173912 277.4668742236025, 804.6287080745341 343.0119456521739, 850.4397795031056 401.50915993788817, 971.6629223602483 387.4134456521739, 950.519350931677 449.4345885093167, 1064.4695652173912 463.00869565217386, 1008.1195652173911 643.5586956521738, 917.2695652173911 659.6586956521738, 908.0695652173912 549.2586956521739, 983.6442795031055 528.3705885093167 )",
		"LINESTRING ( 655.2151321964416 235.9504093665589, 645.5661777376428 253.4939629280114, 634.1628679226986 268.4059834552461, 654.337954518369 284.1951816605533, 665.7412643333131 299.107202187788, 656.9694875525869 317.52793342731314, 601.7072938340115 318.40511110538574, 566.6201867111064 309.63333432465953, 547.3222777935085 317.52793342731314, 541.1820340470002 335.0714869887657, 555.2168768961623 351.73786287214557, 578.0234965260505 361.3868173309445, 580.6550295602684 378.0531932143244, 581.532207238341 398.2282798099948, 592.0583393752125 407.00005659072104, 604.3388268682293 424.5436101521736, 614.8649590051008 439.4556306794082, 617.4964920393187 451.736118172425, 615.7421366831734 466.6481386996596, 607.8475375805198 487.7004029734027, 590.3039840190672 495.59500207605635, 565.7430090330337 497.34935743220154, 550.830988505799 492.96346904183844, 521.8841251294024 480.68298154882166, 499.07750549951413 472.78838244616804, 489.4285510407152 469.27967173387754, 456.09579927395544 471.0340270900228, 453.4642662397375 478.9286261926764, 463.99039837660905 499.1037127883468, 472.7621751573353 508.7526672471457, 485.0426626503521 521.9103324182352, 492.9372617530057 533.3136422331793, 495.5687947872236 538.576708301615, 496.4459724652962 549.9800181165591, 501.709038533732 553.4887288288496, 512.2351706706036 545.594129726196, 515.743881382894 539.4538859796877, 524.5156581636203 542.9625966919781, 523.6384804855477 554.3659065069223, 510.48081531445825 559.628972575358, 503.46339388987724 557.8746172192128, 500.8318608556594 540.3310636577603, 504.3405715679499 524.541865452453, 520.1297697732572 521.0331547401624, 537.6733233347097 536.8223529454698, 530.6559019101287 553.4887288288496, 517.4982367390393 565.7692163218665, 506.97210460216775 569.277927034157, 498.2003278214415 561.3833279315033, 507.8492822802404 541.2082413358329, 517.4982367390393 542.9625966919781, 533.2874349443465 556.1202618630675, 522.761302807475 568.4007493560844, 511.3579929925309 571.9094600683749, 508.726459958313 591.2073689859726, 498.2003278214415 597.347612732481, 472.7621751573353 586.8214805956095, 458.7273323081733 564.8920386437939, 447.32402249322917 544.7169520481234, 430.6576466098492 520.1559770620898, 415.7456260826146 500.85806814449205, 387.6759403842905 487.7004029734027, 363.9921430763296 493.8406467199111, 348.20294487102234 512.2613779594362, 335.04527969993296 516.6472663497993, 314.87019310426257 508.7526672471457, 307.85277167968155 499.1037127883468, 302.5897056112458 488.57758065147533, 292.94075115244686 471.0340270900228, 281.53744133750274 462.2622503092965, 246.4503342145977 460.50789495315126, 233.29266904350828 465.770961021587, 220.1350038724189 503.48960117870996, 222.7665369066368 543.8397743700508, 243.8188011803798 577.1725261368106, 284.16897437172065 591.2073689859726, 306.0984163235363 601.7335011228441, 349.080122549095 619.2770546842967, 394.69336180887154 649.1010957387659, 370.13238682283804 666.6446493002185, 379.7813412816369 681.5566698274531, 411.3597376922515 695.5915126766151, 451.7099108835923 696.4686903546877, 470.1306421231174 689.4512689301067, 497.3231501433689 688.5740912520341, 535.0417903004918 694.7143349985425, 564.8658313549611 724.5383760530118, 561.3571206426706 753.4852394294085, 534.1646126224192 778.046214415442, 494.69161710915097 792.081057264604, 459.60450998624594 792.081057264604, 396.4477171650168 786.8179911961684 )",
		"LINESTRING ( 98.76029805852247 471.39861407224487, 132.68470520765234 514.2079849985279, 157.72414857962914 532.7856365325752, 194.9264255896295 552.0343728591217, 197.44131058055402 568.7784526288414, 256.26647410805407 582.8645232765288, 289.38315727744276 587.7108671549759, 330.5770802442433 581.2490753170464, 362.8860394338908 554.5941839855872, 379.0405190287146 538.4397043907635, 388.73320678560884 507.74619316059835, 419.4267180157739 469.7831661127625, 407.3108583196561 434.24331100415026, 383.8868629071617 410.01159161191464, 355.61652361622015 399.5111798752792, 311.999428710196 393.0493880373497, 250.61240624986573 397.08800793605565, 226.3806868576301 420.5120033485501, 206.9953113438416 457.6673064166447, 202.95669144513568 499.66895336318646, 194.0717276679826 537.6319804110223, 194.9264255896295 552.0343728591217, 199.7257955261709 559.4405278640344, 222.34206695892416 590.9417630739407, 266.7668858446895 620.0198263446234, 292.6140531964075 628.0970661420353, 328.9616322847609 630.5202380812589, 362.8860394338908 621.6352743041058, 427.50395781318576 586.9031431752347, 454.158849144645 542.4783242894694, 463.043812921798 506.93846918085717, 467.082432820504 469.7831661127625, 480.006016496363 438.2819309028562, 505.85318384808096 402.74207579424393, 587.433305801941 375.2794604830436, 624.5886088700355 382.54897630071423, 642.3585364243417 412.4347635511382, 643.4709527477502 439.116190719288, 605.3439380779691 452.77166096041014, 574.5798901784301 466.2701717734732, 531.700351199799 481.0913018291391, 523.623111402387 491.59171356577457, 520.3922154834223 515.8234329580102, 521.1999394631636 540.0551523502459, 530.8926272200578 545.7092202084342, 581.1879051657643 561.1022524449884, 599.9685326993344 539.6722729770875, 620.4120431716209 520.2642150257254, 632.0270408479774 503.31259679536714, 643.4709527477502 439.116190719288, 653.6666721407183 466.55227019379777, 656.8975680596831 497.2457814239629, 644.7817083635653 540.862876329987, 593.0873736601293 582.0567992967876, 551.8934506933286 597.4035549118702, 542.0924468800964 637.2061561170797, 560.9275782471613 680.6088501368376, 577.3059533489567 719.098031626057, 558.4708219818918 740.389919258391, 542.2007629364344 739.5629753463191, 522.8153874226459 758.1406268803665, 479.1982925166218 776.7182784144138, 446.8893333269743 777.5260023941551, 404.88768638043257 754.1020069816606, 393.57955066405594 727.4471156502013, 438.8120935295624 649.9056135950474, 464.65926088128043 641.8283737976354, 542.0924468800964 637.2061561170797, 580.1637899842702 662.0214732911652, 608.4341292752118 686.2531926834008, 656.8975680596831 763.7946947385548, 666.5902558165773 816.296753421732 )",
		"LINESTRING ( 195.66484539313151 434.97628883283056, 229.2424024032942 507.04519168391147, 252.17341694681994 560.2779040170963, 295.57855161849363 579.9330593401183, 354.54401758755984 607.7778627143996, 375.01813771570784 620.0623347912883, 406.1388003104928 606.9588979092737, 439.71635732065545 568.4675520683554, 446.26807576166283 531.6141358376891, 485.57838640770694 456.2693737661045, 501.1387177050994 411.22630948417896, 526.526626664003 390.75218935603095, 578.9403741920618 377.64875247401625, 626.4403328893651 393.2090837714087, 641.1816993816316 422.6918167559418, 646.9144530175131 452.1745497404749, 614.1558608124762 477.5624586993784, 584.6731278279432 481.657282725008, 561.7421132844174 478.3814235045043, 543.7248875716472 484.9331419455117, 529.8024858845066 508.68312129416336, 530.6214506896325 536.5279246684446, 545.362817181899 560.2779040170963, 587.9489870484468 580.7520241452443, 653.4661714585204 592.2175314170071, 699.3282005455719 602.0451090785182, 717.3454262583422 609.4157923246514, 741.9143704121198 643.81231413994, 750.1040184633789 674.1140119295991, 718.983355868594 722.4329354320283, 649.3713474328908 751.9156684165614, 601.0524239304616 756.8294572473169, 505.23354173072903 759.2863516626946, 472.47494952569224 756.010492442191, 429.0698148540185 735.536372314043, 415.9663779720038 681.4846951757323, 424.9749908283889 667.5622934885916, 503.5956121204772 649.5450677758215, 522.4318026383734 644.6312789450659, 516.6990490024918 616.7864755707847, 479.0266679666996 598.7692498580144, 387.3026097925966 611.8726867400292, 343.8974751209229 669.2002230988435, 316.05267174664164 701.1398504987544, 257.08720577757543 699.5019208885026 )",
		"LINESTRING ( 550.167350871396 350.310331871727, 592.2536152006655 404.96932264375374, 598.1767924140437 445.90650202237117, 605.7156815296942 489.107237200183, 601.5087858018727 522.7624030227547, 581.3156863083296 537.907227642912, 566.1708616881724 543.796881661862, 548.5018996313223 545.4796399529906, 532.5156958656007 539.5899859340406, 512.3225963720577 520.2382655860619, 524.1019044099578 494.99689121913303, 540.0881081756794 486.5830997634901, 553.550174504708 483.21758318123295, 563.6467242514796 479.0106874534115, 572.0605157071225 477.32792916228294, 583.8398237450226 466.3900002699471, 591.0281595282754 439.9108744407591, 570.3777574159939 412.54173495383236, 544.2950039035009 405.81070178931805, 508.9570797898005 410.0175975171395, 499.7019091885933 431.0520761562468, 485.39846371400034 475.64517087115433, 482.03294713174313 496.67964951026164, 482.03294713174313 541.2727442251692, 498.86053004302903 569.0382560287908, 519.053629536572 601.010663560234, 540.9294873212436 625.4106587815984, 564.4881033970438 685.1485781166632, 562.8053451059152 732.2658102682636, 524.1019044099578 770.969250964221, 460.1570893470715 780.2244215654283, 419.7708903599855 769.2864926730925, 392.84675770192814 733.9485685593922, 382.7502079551566 694.4037487178705, 390.32262026523523 668.3209952053774, 416.4053737777283 655.700308021913, 437.43985241683566 627.093417072727, 424.8191652333712 597.6451469779768, 403.7846865942639 590.0727346678981, 366.76400418943507 586.707218085641, 327.2191843479133 598.486526123541, 316.2812554555775 633.8244502372413, 282.6260896330058 675.893407515456, 270.84678159510565 721.3278813759277, 254.01919868381984 736.4727059960851, 222.0467911523767 735.6313268505207, 168.19852583626195 699.4520235912562, 168.19852583626195 669.1623743509416, 165.67438839956907 637.1899668194985, 165.67438839956907 615.314109034827, 171.5640424185191 543.796881661862, 174.9295590007763 499.20378694695455, 139.591634887076 467.2313794155114, 109.30198564676145 441.99000504858265, 137.90887659594742 415.9072515360895, 196.8054167854479 406.65208093488235, 263.27436928502703 384.77622315021074, 323.01228862009185 399.0796686248037, 357.50883358822784 431.8934553018111, 369.2881416261279 458.8175879598685, 370.97089991725653 483.21758318123295, 367.6053833349993 512.6658532759832, 365.08124589830646 558.9417062820193, 295.2467768164702 585.0244597945124, 238.87437406366257 591.7554929590267, 216.15713713342666 541.2727442251692 )",
		"LINESTRING ( 792.1521279243493 458.85841310138403, 748.9791978212621 473.6991078243202, 739.5351193612119 500.68218913874966, 721.9961165068327 510.1262675988, 709.8537299153395 510.1262675988, 692.3147270609603 504.72965133591407, 673.4265701408597 493.9364188101423, 650 500, 619.5321667547868 559.3751155970295, 579.4569174196099 540.7835050807104, 550.2887462540068 549.4595899883259, 519.1374699666633 520.5393069629406, 479.47536753184903 499.88196194480827, 459.2311694140793 460.219859509994, 479.47536753184903 440.38880829258693, 479.47536753184903 380.482507740003, 556.0728028590839 392.4637678505197, 619.1190198544241 379.6562139392777, 659.194269189601 441.2151020933122, 650 500, 647.7926428921518 535.7601948475079, 647.7926428921518 558.695813964773, 645.0943347607088 574.8856627534307, 628.9044859720511 587.0280493449239, 618.1112534462793 647.7399823023901, 623.5078697091652 705.7536071284135, 622.1587156434438 730.0383803113999, 610.0163290519505 766.4655400858798, 565.494244883142 796.1469295317521, 543.9077798315984 792.0994673345876, 502.08400379423273 788.0520051374233, 422.48391391666587 766.4655400858798, 383.3584460107432 713.8485315227423, 384.70760007646464 669.3264473539338, 379.3109838135788 626.1535172508466, 369.86690535352847 573.5365086877092, 356.37536469631374 535.7601948475079, 322.64651305327686 477.7465700214847, 276.7752748187468 446.7160265098908, 247.09388537287444 435.922793984119, 213.36503372983765 434.57363991839753, 180.9853361525223 452.11264277277667, 160.74802516670022 512.8245757302429, 170.1921036267505 558.695813964773, 205.2701093355088 624.8043631851251, 252.49050163576032 692.2620664711988, 270.02950449013946 721.9434559170711, 291.615969541683 774.5604644802086, 315.9007427246695 793.4486214003092, 395.5008326022364 827.1774730433459, 427.8805301795518 839.3198596348392, 534.463701371548 847.414784029168, 607.3180209205076 844.7164758977251, 676.1248782723027 816.3842405175742, 720.6469624411112 786.7028510717018, 763.8198925441984 748.9265372315006, 798.8978982529567 716.5468396541852, 843.4199824217652 667.9772932882122, 865.0064474733088 591.0755115420883, 871.7522178019161 545.2042733075583, 893.3386828534597 439.97025618128345, 905.4810694449529 391.40070981531045, 925.718380430775 360.3701663037166 )",
		"LINESTRING ( 471.4485320397647 374.794315485739, 472.1966587531694 399.4824970280944, 478.9297991738118 430.1556922776874, 472.1966587531694 446.61447997259097, 477.05948239030005 471.30266151494635, 496.8848402955248 490.75395606346876, 508.8548677100001 512.8236941089076, 521.1989584811778 528.1602917337042, 546.6352667369379 536.389685581156, 566.1708616881724 543.796881661862, 581.3156863083296 537.907227642912, 599.0041366752674 536.389685581156, 609.8519740196357 522.1752780264665, 619.9516846505992 506.4646170449676, 629.303268568158 489.6317659933617, 624.4404449310275 478.0358019355887, 615.0888610134687 468.68421801802987, 601.6225801721839 467.18796459122046, 594.5153763948392 466.8139012345181, 593.7672496814345 481.4023721459099, 601.2485168154816 492.6242728469805, 602.3707068855886 501.97585676453934, 596.7597565350533 512.8236941089076, 585.5378558339827 516.9383910326335, 563.4681177885437 516.5643276759312, 554.116533870985 515.0680742491218, 537.6577461760813 507.2127437583723, 533.9171126090579 494.12052627378995, 539.5280629595932 482.8986255727193, 551.4980903740685 481.02830878920753, 564.9643712153533 482.52456221601693, 573.193765062805 491.8761461335758, 578.0565886999356 498.98334991092054, 581.4231589102568 497.86115984081346, 583.6675390504709 492.99833620368287, 583.2934756937685 481.4023721459099, 580.675032196852 476.91361186548164, 571.3234482792932 470.18047144483927, 552.9943438008779 466.8139012345181, 530.5505423987366 463.8213943808993, 510.7251844935119 451.47730360972156, 503.99204407286953 434.64445255811563, 508.1067409965954 415.1931580095932, 521.1989584811778 406.21563744873674, 555.9868506544967 403.22313059511794, 577.3084619865309 408.46001758895085, 591.148806184518 429.78162892098504, 606.1113404526121 444.37009983237687, 615.8369877268733 445.4922899024839, 624.0663815743251 442.4997830488651, 628.1810784980511 430.9038189910921, 633.0439021351816 403.97125730852264, 623.6923182176228 387.8865329703214, 611.3482274464451 377.03869562595315, 600.5003901020768 362.07616135785895, 575.4381452030191 343.7470568794436, 529.0542889719272 331.40296610826596, 493.892333441906 329.5326493247542, 456.86006112837293 337.76204317220595, 428.05718266229167 381.90151926308374, 426.56092923548226 416.3153480797003, 435.53844979633874 473.54704165516046, 449.0047306376235 488.5095759232546, 475.5632289634906 514.6940108924194, 514.8398814172377 541.6265725749889, 539.5280629595932 550.230029779143, 566.4606246421627 556.9631701997854, 589.6525527577086 562.948183907023, 621.0738747207063 563.3222472637254, 652.6074156907149 574.2448972794342, 665.5126014969461 591.4518116877425, 650.0263785294686 618.9828747410357, 600.1263267453746 613.8208004185433, 578.6176837349892 617.2621833002049, 568.2935350900043 626.7259862247745, 557.9693864450193 637.9104805901749, 500.32622317718653 645.6535920739136, 485.70034593012446 655.1173949984832, 464.1917029197391 683.5088037721919, 456.44859143600036 727.3864355133779, 459.02962859724664 759.2192271687483, 478.8175801668012 782.4485616199645, 524.4159033488181 798.7951303078573, 569.1538808104197 798.7951303078573, 599.2659810249592 778.1468330178874, 607.8694382291133 748.8950785237633, 598.4056353045437 720.5036697500547, 593.2435609820512 712.7605582663159, 596.6849438637129 691.2519152559306, 568.2935350900043 663.7208522026373, 528.7176319508952 643.9329006330828, 502.04691461801735 639.6311720310057, 477.0968887259703 642.212209192252, 447.84513423184626 656.838086439314, 428.91752838270713 676.6260380088685, 400.52611960899844 692.112260976346, 346.32433922282735 704.1571010621618, 308.4691275245491 692.9726066967614, 297.2846331591487 668.0225808047144, 305.88809036330287 644.7932463534981, 329.11742481451904 609.5190718164662, 342.8829563411657 594.8931945694042, 360.95021646988937 580.2673173223421, 379.8778223190285 551.8759085486334, 379.8778223190285 501.9758567645394, 356.6484878678123 476.165485152077, 327.3967333736882 460.67926218459957, 285.2397930733329 455.51718786210705, 258.5690757404551 459.8189164641841, 237.06043273006972 472.72410227041536, 209.52936967677647 500.2551653237086, 217.2724811605152 561.339711473203, 238.8743740636626 591.7554929590267, 270.6139158262709 611.239763257297, 287.82083023457915 649.9553206759907, 280.07771875084046 692.112260976346, 256.8483842996243 703.2967553417463, 208.66902395636106 698.1346810192539, 173.74784636583726 672.4449405809847, 163.0707007743441 650.8156663964061, 150.1655149681129 609.5190718164662, 152.74655212935912 563.9207486344493, 156.18793501102078 532.9483026994943, 174.9295590007763 499.2037869469546, 210.38971539719188 461.539607905015, 241.3621613321468 427.98612480881377, 248.2449270954701 418.52232188424426, 255.98803857920885 394.4326417126126, 223.2949012034231 371.20330726139645, 181.1379609030678 372.92399870222727, 101.98615462484969 454.65684214169164, 88.22062309820305 575.1052429998497, 120.9137604739888 642.212209192252, 157.0482807314362 692.9726066967614, 204.36729535428398 740.2916213196092, 237.92077845048516 767.8226843729025, 293.8432502774871 793.6330559853649, 236.2000870096543 852.1365649736131, 151.8862064089437 833.208959124474, 90.80166025944929 804.8175503507653, 44.34299135701692 767.8226843729025 )",
		"LINESTRING ( 161.87880870364123 450.84767510297326, 161.87880870364123 488.10066657145876, 179.20578147968095 514.9574743743203, 179.20578147968095 541.814282177182, 158.41341414843328 562.6066495084297, 160.14611142603724 600.7259896157171, 198.2654515333247 636.2462838065986, 242.44923211222604 647.5088161110244, 274.5041317478996 644.9097701946184, 326.48505007601887 630.1818433349847, 348.14376604606855 605.924081448529, 371.5351792937222 547.0123740099939, 371.5351792937222 519.2892175683303, 356.8072524340884 467.30829924021106, 324.75235279841485 445.6495832701614, 296.16284771794926 437.85244552094343, 277.1031776643056 439.58514279854745, 254.57811305545388 448.2486291865673, 233.78574572420618 475.9717856282309, 230.32035116899823 509.75938254150844, 240.71653483462208 544.4133280935879, 258.9098562494638 578.2009250068654, 276.23682902550354 586.8644113948853, 323.0196555208109 589.4634573112912, 343.8120228520586 575.6018790904594, 349.01011468487053 555.6758603980138, 338.61393101924665 529.6854012339542, 327.35139871482085 509.75938254150844, 312.62347185518706 484.6352720162508, 293.5638018015433 469.907345156617, 292.69745316274134 475.9717856282309, 309.1580772999791 505.4276393474985, 316.95521504919697 532.28444715036, 318.6879123268009 552.2104658428058, 307.4253800223751 561.7403008696276, 299.62824227315724 567.8047413412415, 281.4349208583155 569.5374386188456, 269.30603991508764 564.3393467860336, 258.04350761066183 555.6758603980138, 251.11271850024593 540.0815848995779, 246.780975306236 525.3536580399442, 245.048278028632 516.6901716519243, 250.83622260496082 501.7352458225647, 240.71653483462208 485.50162065505276, 242.44923211222604 479.43718018343884, 252.8454157778499 475.10543698942894, 264.10794808227575 471.64004243422096, 273.6377831090976 465.57560196260704, 301.3609395507612 456.9121155745872, 318.6879123268009 462.9765560462011, 336.8812337416427 482.03622609984484, 351.60916060127647 510.6257311803104, 359.40629835049435 534.883493066766, 355.94090379528643 554.8095117592118, 351.60916060127647 571.2701358964496, 343.8120228520586 587.7307600336873, 312.62347185518706 598.1269436993111, 287.4993613299294 605.924081448529, 257.17715897185985 606.7904300873311, 236.38479164061215 593.7952005053012, 213.85972703176046 556.5422090368157, 212.1270297541565 515.8238230131224, 212.1270297541565 498.4968502370826, 215.59242430936445 468.17464787901304, 238.9838375570181 449.1149778253693, 265.8406453598797 426.58991321651763, 296.16284771794926 417.0600781896958, 326.48505007601887 410.99563771808187, 368.0697847385142 421.3918213837057, 395.7929411801778 448.2486291865673, 426.9814921770494 475.9717856282309, 450.37290542470305 483.7689233774488, 485.89319961558454 493.29875840427064, 522.2798424452681 488.96701521026074, 569.9290175793773 481.16987746104286, 608.9147063254668 475.9717856282309, 658.2965787371801 458.64481285219114, 705.9457538712894 417.0600781896958, 700.7476620384774 392.8023163032401, 660.895624653586 361.6137653063686, 595.053128104635 354.6829761959527, 518.81444789006 366.81185713918046, 498.88842919761436 378.9407380824083, 486.7595482543865 403.198499968864, 493.6903373648024 434.3870509657355, 503.22017239162426 448.2486291865673, 521.413493806466 461.24385876859714, 544.8049070541197 474.23908835062696, 562.1318798301594 479.43718018343884, 573.3944121345853 483.7689233774488, 578.5925039673972 491.5660611266667, 584.6569444390111 521.0219148459342, 561.2655311913575 540.94793353838, 514.4827046960502 551.3441172040039, 479.8287591439706 555.6758603980138, 446.9075108694951 573.0028331740535, 423.5160976218414 603.325035532123, 419.1843544278315 644.9097701946184, 434.7786299262673 657.0386511378463, 468.5662268395448 671.7665779974801, 517.0817506124561 678.697367107896, 558.6664852749515 696.8906885227377, 572.5280634957833 723.7474963255993, 567.3299716629714 749.7379554896589, 535.2750720272978 768.7976255433027, 473.7643186723567 768.7976255433027, 419.1843544278315 753.2033500448669, 405.3227762069997 742.8071663792431, 406.1891248458017 711.6186153823716, 417.4516571502275 688.2272021347178, 462.5017863679309 652.7069079438363, 506.68556694683224 640.5780270006086, 571.6617148569813 637.1126324454005, 632.3061195731204 648.3751647498265, 669.5591110416059 664.8357888870642, 696.4159188444675 674.365623913886, 732.802561674151 656.1723024990443, 791.7142691126861 592.9288518664993, 805.5758473335179 553.0768144816078, 806.44219597232 520.1555662071323, 805.5758473335179 489.8333638490627, 792.5806177514881 462.1102074073991 )",
	}
}

func main() {
	options := &opts.Opts{
		Threshold:              50.0,
		MinDist:                20.0,
		RelaxDist:              30.0,
		PlanarSelf:             false,
		AvoidNewSelfIntersects: true,
		GeomRelation:           true,
		DirRelation:            false,
		DistRelation:           false,
	}

	var constraints = []geom.Geometry{}
	for _, wkt := range constraintsWkt {
		constraints = append(constraints, geom.NewPolygonFromWKT(wkt))
	}

	wkt := "LINESTRING ( 161.87880870364123 450.84767510297326, 161.87880870364123 488.10066657145876, 179.20578147968095 514.9574743743203, 179.20578147968095 541.814282177182, 158.41341414843328 562.6066495084297, 160.14611142603724 600.7259896157171, 198.2654515333247 636.2462838065986, 242.44923211222604 647.5088161110244, 274.5041317478996 644.9097701946184, 326.48505007601887 630.1818433349847, 348.14376604606855 605.924081448529, 371.5351792937222 547.0123740099939, 371.5351792937222 519.2892175683303, 356.8072524340884 467.30829924021106, 324.75235279841485 445.6495832701614, 296.16284771794926 437.85244552094343, 277.1031776643056 439.58514279854745, 254.57811305545388 448.2486291865673, 233.78574572420618 475.9717856282309, 230.32035116899823 509.75938254150844, 240.71653483462208 544.4133280935879, 258.9098562494638 578.2009250068654, 276.23682902550354 586.8644113948853, 323.0196555208109 589.4634573112912, 343.8120228520586 575.6018790904594, 349.01011468487053 555.6758603980138, 338.61393101924665 529.6854012339542, 327.35139871482085 509.75938254150844, 312.62347185518706 484.6352720162508, 293.5638018015433 469.907345156617, 292.69745316274134 475.9717856282309, 309.1580772999791 505.4276393474985, 316.95521504919697 532.28444715036, 318.6879123268009 552.2104658428058, 307.4253800223751 561.7403008696276, 299.62824227315724 567.8047413412415, 281.4349208583155 569.5374386188456, 269.30603991508764 564.3393467860336, 258.04350761066183 555.6758603980138, 251.11271850024593 540.0815848995779, 246.780975306236 525.3536580399442, 245.048278028632 516.6901716519243, 250.83622260496082 501.7352458225647, 240.71653483462208 485.50162065505276, 242.44923211222604 479.43718018343884, 252.8454157778499 475.10543698942894, 264.10794808227575 471.64004243422096, 273.6377831090976 465.57560196260704, 301.3609395507612 456.9121155745872, 318.6879123268009 462.9765560462011, 336.8812337416427 482.03622609984484, 351.60916060127647 510.6257311803104, 359.40629835049435 534.883493066766, 355.94090379528643 554.8095117592118, 351.60916060127647 571.2701358964496, 343.8120228520586 587.7307600336873, 312.62347185518706 598.1269436993111, 287.4993613299294 605.924081448529, 257.17715897185985 606.7904300873311, 236.38479164061215 593.7952005053012, 213.85972703176046 556.5422090368157, 212.1270297541565 515.8238230131224, 212.1270297541565 498.4968502370826, 215.59242430936445 468.17464787901304, 238.9838375570181 449.1149778253693, 265.8406453598797 426.58991321651763, 296.16284771794926 417.0600781896958, 326.48505007601887 410.99563771808187, 368.0697847385142 421.3918213837057, 395.7929411801778 448.2486291865673, 426.9814921770494 475.9717856282309, 450.37290542470305 483.7689233774488, 485.89319961558454 493.29875840427064, 522.2798424452681 488.96701521026074, 569.9290175793773 481.16987746104286, 608.9147063254668 475.9717856282309, 658.2965787371801 458.64481285219114, 705.9457538712894 417.0600781896958, 700.7476620384774 392.8023163032401, 660.895624653586 361.6137653063686, 595.053128104635 354.6829761959527, 518.81444789006 366.81185713918046, 498.88842919761436 378.9407380824083, 486.7595482543865 403.198499968864, 493.6903373648024 434.3870509657355, 503.22017239162426 448.2486291865673, 521.413493806466 461.24385876859714, 544.8049070541197 474.23908835062696, 562.1318798301594 479.43718018343884, 573.3944121345853 483.7689233774488, 578.5925039673972 491.5660611266667, 584.6569444390111 521.0219148459342, 561.2655311913575 540.94793353838, 514.4827046960502 551.3441172040039, 479.8287591439706 555.6758603980138, 446.9075108694951 573.0028331740535, 423.5160976218414 603.325035532123, 419.1843544278315 644.9097701946184, 434.7786299262673 657.0386511378463, 468.5662268395448 671.7665779974801, 517.0817506124561 678.697367107896, 558.6664852749515 696.8906885227377, 572.5280634957833 723.7474963255993, 567.3299716629714 749.7379554896589, 535.2750720272978 768.7976255433027, 473.7643186723567 768.7976255433027, 419.1843544278315 753.2033500448669, 405.3227762069997 742.8071663792431, 406.1891248458017 711.6186153823716, 417.4516571502275 688.2272021347178, 462.5017863679309 652.7069079438363, 506.68556694683224 640.5780270006086, 571.6617148569813 637.1126324454005, 632.3061195731204 648.3751647498265, 669.5591110416059 664.8357888870642, 696.4159188444675 674.365623913886, 732.802561674151 656.1723024990443, 791.7142691126861 592.9288518664993, 805.5758473335179 553.0768144816078, 806.44219597232 520.1555662071323, 805.5758473335179 489.8333638490627, 792.5806177514881 462.1102074073991 )"

	var coords = geom.NewLineStringFromWKT(wkt).Coordinates()
	var homo = constdp.NewConstDP(coords, constraints, options, offset.MaxOffset)

	var ptset = homo.Simplify().SimpleSet
	var simplx = []*geom.Point{}

	for _, i := range ptset.Values() {
		simplx = append(simplx, coords[i.(int)])
	}
	fmt.Println(ptset.String())
	fmt.Println(geom.NewLineString(simplx).WKT())
	//ptset = SSet()
}
